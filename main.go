package main

import (
	"flag"
	"fmt"
	"github.com/ordarr/books/service"
	pb "github.com/ordarr/books/v1"
	"github.com/ordarr/data/core"
	"google.golang.org/grpc"
	"log"
	"net"
)

var (
	port = flag.Int("port", 50052, "The server port")
)

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	config, err := core.BuildConfig()
	if err != nil {
		log.Fatalf("failed to build config: %v", err)
	}
	connect := core.Connect(config)
	pb.RegisterBooksServer(s, service.NewServer(core.BookRepository{DB: connect}))
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
