package grpc

import (
	"context"
	"errors"

	pb "github.com/dusk-chancellor/dc-protos/gen/go/sso"
	tkn "github.com/dusk-chancellor/dc-sso/pkg/token"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// token service methods

// validates token
func (s *serverAPI) ValidateToken(ctx context.Context, in *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	token := in.GetToken()
	if token == "" {
		return nil, status.Error(codes.InvalidArgument, "token required")
	}

	valid, err := s.service.ValidateToken(token)
	if err != nil {
		if errors.Is(err, tkn.ErrInvalidToken) {
			return nil, status.Error(codes.InvalidArgument, "invalid token")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &pb.ValidateTokenResponse{
		Valid: valid,
	}, nil
}

// refreshes both access & refresh token based on provided token
func (s *serverAPI) RefreshToken(ctx context.Context, in *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	token := in.GetToken()
	if token == "" {
		return nil, status.Error(codes.InvalidArgument, "token required")
	}

	accessToken, refreshToken, err := s.service.RefreshToken(token)
	if err != nil {
		if errors.Is(err, tkn.ErrInvalidToken) {
			return nil, status.Error(codes.InvalidArgument, "invalid token")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &pb.RefreshTokenResponse{
		AccessToken: accessToken,
		RefreshToken: refreshToken,
	}, nil
}
