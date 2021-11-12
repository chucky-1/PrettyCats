package main

import (
	grpc "CatsCrud/grpc/proto"
	"context"
	"github.com/labstack/gommon/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net"
)

const (
	port = ":10000"
)

type catsCrudServer struct {
	grpc.UnimplementedCatsCrudServer
}

func (s *catsCrudServer) GetAllCats(ctx context.Context, r *grpc.Request) (*grpc.Response, error) {

	return nil, status.Errorf(codes.Unimplemented, "method GetAllCats not implemented")
}
func (s *catsCrudServer) CreateCats(ctx context.Context, r *grpc.Request) (*grpc.Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateCats not implemented")
}
func (s *catsCrudServer) GetCat(ctx context.Context, r *grpc.Request) (*grpc.Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCat not implemented")
}
func (s *catsCrudServer) UpdateCat(ctx context.Context, r *grpc.Request) (*grpc.Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateCat not implemented")
}
func (s *catsCrudServer) DeleteCat(ctx context.Context, r *grpc.Request) (*grpc.Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteCat not implemented")
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
