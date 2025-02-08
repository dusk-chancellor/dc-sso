package service

import (
	"context"
	"errors"

	"github.com/dusk-chancellor/dc-sso/internal/config"
	"github.com/dusk-chancellor/dc-sso/internal/models"
	"github.com/dusk-chancellor/dc-sso/internal/repo"
	"go.uber.org/zap"
)

// errors
var (
	ErrWrongPassword = errors.New("wrong password")
)

// db methods abstraction
type DB interface {
	CreateUser(ctx context.Context, username, email string, password []byte) (id string, err error)
}

// redis methods abstraction for caching

// for get methods
type Getter interface {
	GetUserBy(ctx context.Context, field repo.Field, val any) (models.User, error)
}

// for update methods
type Updater interface {
	UpdateByID(ctx context.Context, id string, field repo.Field, val any) error
	UpdateRole(ctx context.Context, email, wantsRole string) error
}

type Service struct {
	log *zap.SugaredLogger
	db DB
	getter Getter
	updater Updater
	jwt *config.JWT
}

// creates new service instance
func New(logger *zap.SugaredLogger, db DB, getter Getter, updater Updater, jwt *config.JWT) *Service {
	return &Service{
		log: logger,
		db: db,
		getter: getter,
		updater: updater,
		jwt: jwt,
	}
}
