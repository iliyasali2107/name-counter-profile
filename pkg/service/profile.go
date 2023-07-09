package service

import (
	"context"
	"net/http"

	"name-counter-url/pkg/db"
	"name-counter-url/pkg/pb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type service struct {
	S db.Storage
	pb.UnimplementedURLServiceServer
}

type Service interface {
	pb.URLServiceServer
	AddURL(context.Context, *pb.AddURLRequest) (*pb.AddURLResponse, error)
	GetURL(context.Context, *pb.GetURLRequest) (*pb.GetURLResponse, error)
}

func NewService(s db.Storage) Service {
	return &service{
		S: s,
	}
}

func (srv *service) AddURL(ctx context.Context, req *pb.AddURLRequest) (*pb.AddURLResponse, error) {
	_, err := srv.S.InsertURL(req.Id, req.Url)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "something unexpected happened")
	}

	return &pb.AddURLResponse{
		Status: http.StatusOK,
	}, nil
}

func (srv *service) GetURL(ctx context.Context, req *pb.GetURLRequest) (*pb.GetURLResponse, error) {
	url, err := srv.S.GetURL(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "NotFound")
	}

	return &pb.GetURLResponse{
		Status: http.StatusOK,
		Url:    url.URL,
	}, nil
}

func (srv *service) SetActiveURL(ctx context.Context, req *pb.SetActiveUrlRequest) (*pb.SetActiveUrlResponse, error) {
	activeURL, err := srv.S.GetActiveURL(req.UrlID)
	if err != nil {
		return nil, err
	}

	if activeURL.ID == req.UrlID {
		return &pb.SetActiveUrlResponse{
			Status: http.StatusForbidden,
		}, nil
	}

	_, err = srv.S.SetActive(req.UrlID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to set active url: %v", err)
	}

	_, err = srv.S.SetNotActive(activeURL.ID)
	if err != nil {
		return nil, err
	}

	return &pb.SetActiveUrlResponse{
		Status: http.StatusOK,
	}, nil
}
