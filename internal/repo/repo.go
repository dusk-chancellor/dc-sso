package repo

import (
	"context"
	"errors"

	"github.com/dusk-chancellor/dc-sso/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/redis/go-redis/v9"
)

// init DB struct & abstractions & global vars

// errors
var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotFound 	 = errors.New("user not found")
)

// Field type represents a string enum which is as in db fields;
// used as search field
type Field string

// constant representations of db fields
const (
	ID 		 Field = "id"
	Email 	 Field = "email"
	Username Field = "username"
	Password Field = "password"
)

// interface abstraction of pgxpool.Pool methods
type DB interface {
	Exec(context.Context, string, ...any) (pgconn.CommandTag, error)
	Query(context.Context, string, ...any) (pgx.Rows, error)
	QueryRow(context.Context, string, ...any) pgx.Row
}

type DBLayer interface {
	GetUserBy(ctx context.Context, field Field, val any) (models.User, error)
	UpdateByID(ctx context.Context, id string, field Field, val any) error
	UpdateRole(ctx context.Context, email, wantsRole string) error
}

type Repo struct {
	db DB
}

type Rdb struct {
	rdb *redis.Client
	db DBLayer
}


func NewDB(db DB) *Repo {
	return &Repo{
		db: db,
	}
}

func NewRdb(rdb *redis.Client, db DBLayer) *Rdb {
	return &Rdb{
		rdb: rdb,
		db: db,
	}
}

// with transactions
func (r *Repo) WithTx(tx pgx.Tx) *Repo {
	return &Repo{
		db: tx,
	}
}
