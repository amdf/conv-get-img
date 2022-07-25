package main

import (
	"log"
	"net"
	"time"

	"github.com/amdf/conv-get-img/internal/config"
	"github.com/amdf/conv-get-img/internal/server"
	pb "github.com/amdf/conv-get-img/svc"
	"google.golang.org/grpc"

	opentracing "github.com/opentracing/opentracing-go"
	jaeger "github.com/uber/jaeger-client-go"
	jcfg "github.com/uber/jaeger-client-go/config"
)

func main() {

	if errload := config.Load(); errload != nil {
		log.Fatal("unable to load configs/config.toml: ", errload.Error())
	}

	newjcfg := jcfg.Configuration{
		Sampler: &jcfg.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &jcfg.ReporterConfig{
			LogSpans:            true,
			BufferFlushInterval: 1 * time.Second,
		},
	}
	tracer, closer, jerr := newjcfg.New(
		"conv-get-img",
		jcfg.Logger(jaeger.StdLogger),
	)
	if jerr != nil {
		log.Fatalln("fail to create tracing", jerr)
	}
	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()

	srv, err := server.NewServer()
	if err != nil {
		log.Fatalln("fail to create server", err)
	}

	lis, err := net.Listen("tcp", config.GetServerAddress())
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
