package server

import (
	"context"

	"github.com/FreakyGranny/anti-brute-force/internal/app"
	"github.com/golang/protobuf/ptypes/empty"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

//go:generate protoc AbfService.proto --go_out=plugins=grpc:. -I ../../api

// Service grpc events service.
type Service struct {
	app app.Application
}

// New returns grpc service.
func New(a app.Application) *Service {
	return &Service{app: a}
}

// Auth check possibility of authentication.
func (s *Service) Auth(ctx context.Context, req *AuthRequest) (*AuthResponse, error) {
	ok, err := s.app.CheckRate(ctx, req.GetLogin(), req.GetPassword(), req.GetIp())
	switch err {
	case nil:
		return &AuthResponse{Ok: ok}, nil
	case app.ErrInvalidArgument:
		return nil, status.Error(codes.InvalidArgument, "ip/mask is not valid")
	default:
		return nil, status.Error(codes.Internal, "unable to check rate")
	}
}

// DropStat drops buckets for given login, ip.
func (s *Service) DropStat(ctx context.Context, req *DropStatRequest) (*empty.Empty, error) {
	err := s.app.DropStat(ctx, req.GetLogin(), req.GetIp())
	if err != nil {
		return nil, status.Error(codes.Internal, "unable to drop stat")
	}

	return &empty.Empty{}, nil
}

// AddToWhiteList adds ip to white list.
func (s *Service) AddToWhiteList(ctx context.Context, req *AddSubnetRequest) (*AddSubnetResponse, error) {
	err := s.app.AddToWhiteList(ctx, req.GetIp(), req.GetMask())
	switch err {
	case nil:
		return &AddSubnetResponse{Ip: req.GetIp(), Mask: req.GetMask()}, nil
	case app.ErrInvalidArgument:
		return nil, status.Error(codes.InvalidArgument, "ip/mask is not valid")
	default:
		return nil, status.Error(codes.Internal, "unable to add to white list")
	}
}

// AddToBlackList adds ip to black list.
func (s *Service) AddToBlackList(ctx context.Context, req *AddSubnetRequest) (*AddSubnetResponse, error) {
	err := s.app.AddToBlackList(ctx, req.GetIp(), req.GetMask())
	switch err {
	case nil:
		return &AddSubnetResponse{Ip: req.GetIp(), Mask: req.GetMask()}, nil
	case app.ErrInvalidArgument:
		return nil, status.Error(codes.InvalidArgument, "ip/mask is not valid")
	default:
		return nil, status.Error(codes.Internal, "unable to add to black list")
	}
}

// RemoveFromWhiteList removes ip from white list.
func (s *Service) RemoveFromWhiteList(ctx context.Context, req *RemoveSubnetRequest) (*empty.Empty, error) {
	err := s.app.RemoveFromWhiteList(ctx, req.GetIp(), req.GetMask())
	switch err {
	case nil:
		return &empty.Empty{}, nil
	case app.ErrInvalidArgument:
		return nil, status.Error(codes.InvalidArgument, "ip/mask is not valid")
	default:
		return nil, status.Error(codes.Internal, "unable to remove from white list")
	}
}

// RemoveFromBlackList removes ip from black list.
func (s *Service) RemoveFromBlackList(ctx context.Context, req *RemoveSubnetRequest) (*empty.Empty, error) {
	err := s.app.RemoveFromBlackList(ctx, req.GetIp(), req.GetMask())
	switch err {
	case nil:
		return &empty.Empty{}, nil
	case app.ErrInvalidArgument:
		return nil, status.Error(codes.InvalidArgument, "ip/mask is not valid")
	default:
		return nil, status.Error(codes.Internal, "unable to remove from black list")
	}
}
