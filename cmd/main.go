package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"

	"github.com/Shopify/sarama"
	"github.com/amdf/conv-get-img/internal/producer"
	pb "github.com/amdf/conv-get-img/svc"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

func (srv ConvGetImageServer) Convert(context.Context, *pb.ConvertRequest) (resp *pb.ConvertResponse, err error) {
	resp = &pb.ConvertResponse{ConvId: "some ID here"}
	return
}

func MakeConvGetImageServer() (srv *ConvGetImageServer, err error) {
	srv = &ConvGetImageServer{}

	srv.syncProducer, err = producer.NewSync()
	if nil != err {
		return
	}
	srv.asyncProducer, err = producer.NewAsync()
	return
}

//TODO: implement Image

func (srv ConvGetImageServer) RunGateway() {
	//TODO: handle errors another way

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Note: Make sure the gRPC server is running properly and accessible
	mux := runtime.NewServeMux()

	conn, err := grpc.DialContext(
		ctx,
		SVCRPCADDR,
		//grpc.WithInsecure(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln("fail to dial rpc", err)
	}

	err = pb.RegisterConvGetImageHandler(ctx, mux, conn)
	if err != nil {
		log.Fatalln("fail to register http", err)
	}
	//TODO: move this to config
	allowedOrigins := []string{
		"http://localhost:8080",
		"http://127.0.0.1:8080"}

	log.Println("starting http server at ", SVCHTTPADDR)
	gatewayServer := &http.Server{
		Addr:    SVCHTTPADDR,
		Handler: cors(mux, allowedOrigins),
	}
	// Start HTTP server (and proxy calls to gRPC server endpoint)
	if err = gatewayServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalln("fail to serve http", err)
	}
}

func cors(h http.Handler, allowedOrigins []string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		providedOrigin := r.Header.Get("Origin")
		matches := false
		for _, allowedOrigin := range allowedOrigins {
			if providedOrigin == allowedOrigin {
				matches = true
				break
			}
		}

		if matches {
			w.Header().Set("Access-Control-Allow-Origin", providedOrigin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, ResponseType, grpc-metadata-log-level")
		}
		if r.Method == "OPTIONS" {
			return
		}
		h.ServeHTTP(w, r)
	})
}

func main() {

	server, err := MakeConvGetImageServer()
	if err != nil {
		log.Fatalln("fail to create server", err)
	}

	lis, err := net.Listen("tcp", SVCRPCADDR)
	if err != nil {
		log.Fatalln("fail to listen", err)
	}

	s := grpc.NewServer(
	//grpc.InTapHandle(server.RateLimiter),
	)
	pb.RegisterConvGetImageServer(s, server)

	log.Println("starting server at ", lis.Addr())

	go server.RunGateway()

	err = s.Serve(lis)
	if err != nil {
		log.Fatalln("server failed", err)
	}

}
