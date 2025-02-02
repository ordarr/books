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
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"net"
)

type booksServer struct {
	pb.UnimplementedBooksServer
	repo core.BookRepository
}

type BooksResult []*pb.Book

func (s *booksServer) GetBooks(_ context.Context, _ *emptypb.Empty) (*pb.BooksResponse, error) {
	result := &BooksResult{}
	if copier.Copy(result, s.repo.GetAll()) != nil {
		return nil, status.Error(codes.Unknown, "unknown error")
	}
	return &pb.BooksResponse{Content: *result}, nil
}

func (s *booksServer) GetBookByTitle(_ context.Context, request *pb.ValueRequest) (*pb.BookResponse, error) {
	result := &pb.Book{}
	book := s.repo.GetByTitle(request.Value)
	if book.ID == "" || copier.Copy(&result, book) != nil {
		return nil, status.Error(codes.NotFound, "book not found")
	}
	return &pb.BookResponse{Content: result}, nil
}

func (s *booksServer) GetBookById(_ context.Context, request *pb.ValueRequest) (*pb.BookResponse, error) {
	result := &pb.Book{}
	book := s.repo.GetById(request.Value)
	if book.ID == "" || copier.Copy(&result, book) != nil {
		return nil, status.Error(codes.NotFound, "book not found")
	}
	return &pb.BookResponse{Content: result}, nil
}

func NewServer(repository core.BookRepository) pb.BooksServer {
	t := &booksServer{
		repo: repository,
	}
	return t
}

func Server(repository core.BookRepository) (pb.BooksClient, func()) {
	buffer := 101024 * 1024
	lis := bufconn.Listen(buffer)

	baseServer := grpc.NewServer()
	pb.RegisterBooksServer(baseServer, NewServer(repository))
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
