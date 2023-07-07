package service

import (
	"context"

	"name-counter-profile/pkg/db"
	"name-counter-profile/pkg/pb"
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
	panic("TODO")
}
