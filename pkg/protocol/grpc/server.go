package grpc

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"

	"google.golang.org/grpc"

	"github.com/metehanakbaba/go-grpc-http-rest-microservice-tutorial/pkg/api/api"
)

// ToDo gRPC Server
func RunServer(ctx context.Context, v1API api.ToDoServiceServer, port string) error {
	listen, err := net.Listen("tcp", ":"+port) // tcp dinleme
	if err != nil {
		return err
	}

	// gRPC Servisini reg edelim
	server := grpc.NewServer()
	api.RegisterToDoServiceServer(server, v1API)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			// Processlist'de Interrupt (^C) saglanirsa devam eden gRPC serverini gracefully stop edelim ve context'i durduralim
			log.Println("gRPC server kapatiliyor...")

			server.GracefulStop()

			<-ctx.Done()
		}
	}()

	log.Println("gRPC Server baslatiliyor...")

	return server.Serve(listen)
}
