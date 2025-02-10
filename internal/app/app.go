package app

import (
	"context"

	"github.com/dusk-chancellor/dc-sso/internal/app/server"
	"github.com/dusk-chancellor/dc-sso/internal/config"
	"github.com/dusk-chancellor/dc-sso/internal/database/postgres"
	"github.com/dusk-chancellor/dc-sso/internal/database/redis"
	"github.com/dusk-chancellor/dc-sso/internal/repo"
	"github.com/dusk-chancellor/dc-sso/internal/service"
	"go.uber.org/zap"
)

// app init

type App struct {
	GRPCServer *server.Server
}

func New(ctx context.Context, log *zap.Logger, cfg *config.Config) *App {
	pool, err := postgres.ConnectDB(ctx, &cfg.Db)
	if err != nil {
		log.DPanic("failed to connect db", zap.Error(err))
	}

	redisClient, err := redis.NewClient(ctx, &cfg.Redis)
	if err != nil {
		log.DPanic("failed to init redis client", zap.Error(err))
	}

	dbRepo := repo.NewDB(pool)
	rdb := repo.NewRdb(redisClient, dbRepo)

	srvc := service.New(log, dbRepo, rdb, rdb, &cfg.Jwt)

	server := server.New(log, *srvc, cfg.GrpcServer.Port)

	return &App{
		GRPCServer: server,
	}
}
