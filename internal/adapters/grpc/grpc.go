package grpc

import (
	"context"

	"github.com/dusk-chancellor/dc-protos/gen/go/sso"
	"github.com/dusk-chancellor/dc-sso/internal/models"
	"google.golang.org/grpc"
)

// transport layer of the app
// grpc methods implemented, all of the check incoming args for emptiness
// desc: https://github.com/dusk-chancellor/dc-protos/

// service logic methods abstraction
type Service interface {
	// Auth
	Register(ctx context.Context, username, email, password string) (id, accessToken, refreshToken string, err error)
	Login(ctx context.Context, username, email, password string) (id, accessToken, refreshToken string, err error)
	Logout(ctx context.Context, token string) (success bool, err error)

	// User
	GetUser(ctx context.Context, id string) (models.User, error)
	UpdateUser(ctx context.Context, id, username, email string) (success bool, err error)
	SetRole(ctx context.Context, email, wantsRole string) (success bool, err error)
	ChangePassword(ctx context.Context, id, oldPassword, newPassword string) (success bool, err error)

	// Token
	ValidateToken(token string) (valid bool, err error)
	RefreshToken(token string) (accessToken, refreshToken string, err error)
}

type serverAPI struct {
	sso.UnimplementedAuthServiceServer
	sso.UnimplementedUserServiceServer
	sso.UnimplementedTokenServiceServer
	service Service
}

// register grpc handlers
func RegisterGrpc(grpcServer *grpc.Server, srv Service) {
	sso.RegisterAuthServiceServer(grpcServer, &serverAPI{
		service: srv,
	})

	sso.RegisterUserServiceServer(grpcServer, &serverAPI{
		service: srv,
	})

	sso.RegisterTokenServiceServer(grpcServer, &serverAPI{
		service: srv,
	})
}
