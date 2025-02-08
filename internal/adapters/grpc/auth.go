package grpc

import (
	"context"
	"errors"

	pb "github.com/dusk-chancellor/dc-protos/gen/go/sso"
	"github.com/dusk-chancellor/dc-sso/internal/repo"
	"github.com/dusk-chancellor/dc-sso/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// auth service methods

// registers user
func (s *serverAPI) Register(ctx context.Context, in *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	uname, email, password := in.GetUsername(), in.GetEmail(), in.GetPassword()
	if uname == "" || email == "" || password == "" {
		return nil, status.Error(codes.InvalidArgument, "username, email and password are required")
	}

	id, accessToken, refreshToken, err := s.service.Register(ctx, uname, email, password)
	if err != nil {
		if errors.Is(err, repo.ErrUserAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &pb.RegisterResponse{
		Id: id,
		AccessToken: accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// logins user
func (s *serverAPI) Login(ctx context.Context, in *pb.LoginRequest) (*pb.LoginResponse, error) {
	uname, email, password := in.GetUsername(), in.GetEmail(), in.GetPassword()
	if uname == "" || email == "" || password == "" {
		return nil, status.Error(codes.InvalidArgument, "username or email and password are required")
	}

	id, accessToken, refreshToken, err := s.service.Login(ctx, uname, email ,password)
	if err != nil {
		if errors.Is(err, repo.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}

		if errors.Is(err, service.ErrWrongPassword) {
			return nil, status.Error(codes.InvalidArgument, "passwords don't match")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &pb.LoginResponse{
		Id: id,
		AccessToken: accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// unimplemented: logouts user
func (s *serverAPI) Logout(ctx context.Context, in *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	return s.UnimplementedAuthServiceServer.Logout(ctx, in)
}
