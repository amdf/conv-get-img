package server

import (
	"context"
	"encoding/json"
	"log"

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

func (srv ConvGetImageServer) Image(ctx context.Context, req *pb.ImageRequest) (body *httpbody.HttpBody, err error) {
	log.Println("get Image with id = ", req.ConvId)
	body, err = nil, status.Errorf(codes.Unimplemented, "method Image not implemented")
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
