package server

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/Shopify/sarama"
	"github.com/opentracing/opentracing-go"

	"github.com/amdf/conv-get-img/internal/config"
	"github.com/amdf/conv-get-img/internal/producer"
	pb "github.com/amdf/conv-get-img/svc"
	"google.golang.org/genproto/googleapis/api/httpbody"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ConvGetImageServer struct {
	pb.UnimplementedConvGetImageServer
	syncProducer  sarama.SyncProducer
	asyncProducer sarama.AsyncProducer
}

func (srv ConvGetImageServer) SaveText(ctx opentracing.SpanContext, text string) {
	tr := opentracing.GlobalTracer().StartSpan(
		"SaveText", opentracing.ChildOf(ctx))
	defer tr.Finish()
	prepMsg := producer.PrepareMessage("texts", []byte(text))
	srv.asyncProducer.Input() <- prepMsg
}

func (srv ConvGetImageServer) Convert(ctx context.Context, req *pb.ConvertRequest) (resp *pb.ConvertResponse, err error) {
	tr := opentracing.GlobalTracer().StartSpan("Convert")
	defer tr.Finish()

	go srv.SaveText(tr.Context(), req.InputText)

	rqdata := ConvertRequestData{
		InputText: req.InputText,
		FontSize:  req.FontSize,
		FontFile:  req.FontFile,
		FontStyle: FontStyles(req.FontStyle),
	}

	rq := ConvertRequest{ConvertRequestData: rqdata}
	rq.ConvID = rqdata.UniqueID()

	buf, err1 := json.Marshal(rq)
	if err1 != nil {
		err = err1
		log.Println("fail to encode: ", err.Error())
		return
	}
	prepMsg := producer.PrepareMessage("convert_requests", buf)
	part, offset, err2 := srv.syncProducer.SendMessage(prepMsg)
	if err2 != nil {
		err = err2
		log.Println("fail to send: ", err.Error())
		return
	}

	log.Println("convert Image with id = ", rq.ConvID, "send to partition = ", part, "offset = ", offset)

	resp = &pb.ConvertResponse{ConvId: rq.ConvID}

	return
}

func (srv ConvGetImageServer) Image(ctx context.Context, req *pb.ImageRequest) (body *httpbody.HttpBody, err error) {
	tr := opentracing.GlobalTracer().StartSpan("Image")
	defer tr.Finish()

	var timeout time.Duration
	timeout, err = time.ParseDuration(config.Get().Storage.Timeout)
	if err != nil {
		return
	}

	ctx2, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	path := config.Get().Storage.Path + "/" + req.ConvId + ".png"

	found := make(chan struct{}, 1)
	go func() {
		stamp := time.Now()
		for time.Since(stamp) < timeout {
			fst, errf := os.Stat(path)
			if nil == errf && fst.Size() > 0 {
				found <- struct{}{}
			}
		}
	}()

	select {
	case <-ctx2.Done():
		err = status.Errorf(codes.NotFound, "image not found")
		log.Println("get Image ", path, "not found")
		return
	case <-found:
	}

	imgdata, err := os.ReadFile(path)
	if err != nil {
		log.Println("error load Image with id = ", req.ConvId)
		err = status.Errorf(codes.DataLoss, "load image error")
	} else {
		body = &httpbody.HttpBody{
			ContentType: "image/png",
			Data:        imgdata,
		}
		log.Println("get Image with id = ", req.ConvId)
	}
	return
}

func NewServer() (srv *ConvGetImageServer, err error) {
	srv = &ConvGetImageServer{}

	srv.syncProducer, err = producer.NewSync()
	if nil != err {
		return
	}
	srv.asyncProducer, err = producer.NewAsync()

	if nil != err {
		return
	}

	go func() {
		for x := range srv.asyncProducer.Successes() {
			log.Println("ok write to ", x.Topic, x.Partition, x.Offset)
		}
	}()
	go func() {
		for x := range srv.asyncProducer.Errors() {
			log.Println("error write to ", x.Msg.Topic, "-", x.Err.Error())
		}
	}()
	return
}
