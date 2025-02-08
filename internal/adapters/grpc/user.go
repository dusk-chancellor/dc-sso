package grpc

import (
	"context"
	"errors"

	pb "github.com/dusk-chancellor/dc-protos/gen/go/sso"
	"github.com/dusk-chancellor/dc-sso/internal/dto"
	"github.com/dusk-chancellor/dc-sso/internal/repo"
	"github.com/dusk-chancellor/dc-sso/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// user service methods

// get user info
func (s *serverAPI) GetUser(ctx context.Context, in *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	id := in.GetId()
	if id == "" {
		return nil, status.Error(codes.InvalidArgument, "id required")
	}

	user, err := s.service.GetUser(ctx, id)
	if err != nil {
		if errors.Is(err, repo.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	out := dto.ToPbUser(&user)

	return &pb.GetUserResponse{
		User: out,
	}, nil
}

// updates user's name or email
func (s *serverAPI) UpdateUser(ctx context.Context, in *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	id, uname, email := in.GetId(), in.GetUsername(), in.GetEmail()
	if id == "" {
		return nil, status.Error(codes.InvalidArgument, "id required")
	}

	success, err := s.service.UpdateUser(ctx, id, uname, email)
	if err != nil {
		if errors.Is(err, repo.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &pb.UpdateUserResponse{
		Success: success,
	}, nil
}

// sets new role to user
func (s *serverAPI) SetRole(ctx context.Context, in *pb.SetRoleRequest) (*pb.SetRoleResponse, error) {
	email, wantsRole := in.GetEmail(), in.GetWantsRole()
	if email == "" || wantsRole.String() == "" {
		return nil, status.Error(codes.InvalidArgument, "email and wanted role required")
	}

	success, err := s.service.SetRole(ctx, email, wantsRole.String())
	if err != nil {
		if errors.Is(err, repo.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &pb.SetRoleResponse{
		Success: success,
	}, nil
}

// changes old password to new
func (s *serverAPI) ChangePassword(ctx context.Context, in *pb.ChangePasswordRequest) (*pb.ChangePasswordResponse, error) {
	id, oldPass, newPass := in.GetId(), in.GetOldPassword(), in.GetNewPassword()
	if id == "" || oldPass == "" || newPass == "" {
		return nil, status.Error(codes.InvalidArgument, "id, old password and new password required")
	}

	success, err := s.service.ChangePassword(ctx, id, oldPass, newPass)
	if err != nil {
		if errors.Is(err, repo.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}

		if errors.Is(err, service.ErrWrongPassword) {
			return nil, status.Error(codes.InvalidArgument, "passwords don't match")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &pb.ChangePasswordResponse{
		Success: success,
	}, nil
}
