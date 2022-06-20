package server

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/Shopify/sarama"
	"github.com/amdf/conv-get-img/internal/producer"
	pb "github.com/amdf/conv-get-img/svc"
	"google.golang.org/genproto/googleapis/api/httpbody"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	SVCRPCADDR  = "0.0.0.0:50052"
	SVCHTTPADDR = "0.0.0.0:8082"
)

type ConvGetImageServer struct {
	pb.UnimplementedConvGetImageServer
	syncProducer  sarama.SyncProducer
	asyncProducer sarama.AsyncProducer
}

func (srv ConvGetImageServer) Convert(ctx context.Context, req *pb.ConvertRequest) (resp *pb.ConvertResponse, err error) {
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

var GETIMGTIMEOUT = 3 * time.Second //TODO: move to config

func (srv ConvGetImageServer) Image(ctx context.Context, req *pb.ImageRequest) (body *httpbody.HttpBody, err error) {
	ctx2, cancel := context.WithTimeout(ctx, GETIMGTIMEOUT)
	defer cancel()

	path := "img/" + req.ConvId + ".png"

	found := make(chan struct{}, 1)
	go func(s string) {
		stamp := time.Now()
		for time.Since(stamp) < GETIMGTIMEOUT {
			fst, errf := os.Stat(s)
			if nil == errf && fst.Size() > 0 {
				found <- struct{}{}
			} else {

				fst, errf := os.Stat("../" + s)
				if nil == errf && fst.Size() > 0 {
					found <- struct{}{}
				}
			}
		}
	}(path)

	select {
	case <-ctx2.Done():
		body, err = nil, status.Errorf(codes.NotFound, "image not found")
		log.Println("get Image ", path, "not found")
	case <-found:
		body = &httpbody.HttpBody{ContentType: "image/png"}
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
	return
}
