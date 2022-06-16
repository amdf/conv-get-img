package server

import (
	"context"
	"errors"
	"log"
	"net/http"

	pb "github.com/amdf/conv-get-img/svc"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var allowedOrigins = []string{
	"http://localhost:8080",
	"http://127.0.0.1:8080",
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
