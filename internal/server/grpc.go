package server

import (
	"context"
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
	resp = &pb.ConvertResponse{ConvId: "some ID here"}
	log.Println("convert Image with len = ", len(req.InputText), "font size", req.FontSize)
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
