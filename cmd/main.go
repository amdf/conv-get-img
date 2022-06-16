package main

import (
	"log"
	"net"

	"github.com/amdf/conv-get-img/internal/server"
	pb "github.com/amdf/conv-get-img/svc"
	"google.golang.org/grpc"
)

func main() {

	srv, err := server.NewServer()
	if err != nil {
		log.Fatalln("fail to create server", err)
	}

	lis, err := net.Listen("tcp", server.SVCRPCADDR)
	if err != nil {
		log.Fatalln("fail to listen", err)
	}

	s := grpc.NewServer(
	//grpc.InTapHandle(server.RateLimiter),
	)
	pb.RegisterConvGetImageServer(s, srv)

	log.Println("starting server at ", lis.Addr())

	go srv.RunGateway()

	err = s.Serve(lis)
	if err != nil {
		log.Fatalln("server failed", err)
	}

}
