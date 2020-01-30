package rest

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/metehanakbaba/go-grpc-http-rest-microservice-tutorial/pkg/api/api"
	"google.golang.org/grpc"
)

// HTTP/REST gateway
func RunServer(ctx context.Context, grpcPort, httpPort string) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	if err := api.RegisterToDoServiceHandlerFromEndpoint(ctx, mux, "localhost:"+grpcPort, opts); err != nil {
		return err
	}

	srv := &http.Server{
		Addr:    ":" + httpPort,
		Handler: mux,
	}

	// Processlist'de Interrupt (^C) saglanirsa devam eden HTTP sunucusu kapatalim
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {

		}

		_, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		_ = srv.Shutdown(ctx)
	}()

	log.Println("HTTP/REST baslatiliyor...")
	return srv.ListenAndServe()
}
