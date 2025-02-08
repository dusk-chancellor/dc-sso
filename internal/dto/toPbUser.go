package dto

import (
	"github.com/dusk-chancellor/dc-sso/internal/models"

	pb "github.com/dusk-chancellor/dc-protos/gen/go/sso"
)

func ToPbUser(in *models.User) *pb.User {
	var role pb.Role
	if in.Role == pb.Role_ADMIN.String() {
		role = pb.Role_ADMIN
	} else {
		role = pb.Role_USER
	}

	return &pb.User{
		Id: in.ID.String(),
		Username: in.Username,
		Email: in.Email,
		Role: role,
	}
}
