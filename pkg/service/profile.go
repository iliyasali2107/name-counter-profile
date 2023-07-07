package service

import (
	"context"
	"net/http"

	"name-counter-profile/pkg/db"
	"name-counter-profile/pkg/pb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type service struct {
	S db.Storage
	pb.UnimplementedProfileServiceServer
}

type Service interface {
	pb.ProfileServiceServer
	AddURL(context.Context, *pb.AddURLRequest) (*pb.AddURLResponse, error)
	GetURL(context.Context, *pb.GetURLRequest) (*pb.GetURLResponse, error)
	SetInterval(context.Context, *pb.SetIntervalRequest) (*pb.SetIntervalResponse, error)
}

func NewService(s db.Storage) Service {
	return &service{
		S: s,
	}
}

func (srv *service) AddURL(ctx context.Context, req *pb.AddURLRequest) (*pb.AddURLResponse, error) {
	_, err := srv.S.InsertURL(int(req.Id), req.Url)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "something unexpected happened")
	}

	return &pb.AddURLResponse{
		Status: http.StatusOK,
	}, nil
}

func (srv *service) GetURL(ctx context.Context, req *pb.GetURLRequest) (*pb.GetURLResponse, error) {
	url, err := srv.S.GetURL(int(req.Id))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "NotFound")
	}

	return &pb.GetURLResponse{
		Status: http.StatusOK,
		Url:    url,
	}, nil
}
