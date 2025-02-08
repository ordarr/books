package service

import (
	"context"
	"github.com/jinzhu/copier"
	pb "github.com/ordarr/books/v1"
	"github.com/ordarr/data/core"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
	"log"
	"net"
)

type booksServer struct {
	pb.UnimplementedBooksServer
	repo core.IBookRepository
}

type BooksResult []*pb.Book

func (s *booksServer) GetBooks(_ context.Context, request *pb.GetBooksRequest) (*pb.BooksResponse, error) {
	result := &BooksResult{}
	var books []*core.Book
	var err error

	if len(request.Ids) == 0 && len(request.Names) == 0 {
		books, err = s.repo.GetAll()
	} else if request.Ids != nil {
		books, err = s.repo.GetByID(request.Ids)
	} else {
		books, err = s.repo.GetByName(request.Names)
	}
	if err != nil {
		return nil, err
	}

	if copier.Copy(&result, books) != nil {
		return nil, status.Error(codes.Unknown, "unknown error")
	}
	return &pb.BooksResponse{Content: *result}, nil
}

func (s *booksServer) CreateBook(_ context.Context, request *pb.CreateBookRequest) (*pb.BookResponse, error) {
	result := &pb.Book{}
	created := &core.Book{
		Name: request.Name,
	}
	created, err := s.repo.Create(created)
	if err != nil {
		return nil, err
	}
	if copier.Copy(&result, created) != nil {
		return nil, status.Error(codes.Unknown, "unknown error")
	}
	return &pb.BookResponse{Content: result}, nil
}

func Server(repository core.IBookRepository) (*grpc.Server, error) {
	baseServer := grpc.NewServer()
	pb.RegisterBooksServer(baseServer, &booksServer{repo: repository})
	return baseServer, nil
}

func CreateClient(repository core.IBookRepository) (pb.BooksClient, func()) {
	buffer := 101024 * 1024
	lis := bufconn.Listen(buffer)

	baseServer := grpc.NewServer()

	pb.RegisterBooksServer(baseServer, &booksServer{repo: repository})

	go func() {
		if err := baseServer.Serve(lis); err != nil {
			log.Printf("error serving server: %v", err)
		}
	}()

	conn, err := grpc.NewClient("localhost:8080",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return lis.Dial()
		}), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("error connecting to server: %v", err)
	}

	closer := func() {
		err := lis.Close()
		if err != nil {
			log.Printf("error closing listener: %v", err)
		}
		baseServer.Stop()
	}

	client := pb.NewBooksClient(conn)

	return client, closer
}
