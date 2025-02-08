package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/dusk-chancellor/dc-sso/internal/app"
	"github.com/dusk-chancellor/dc-sso/internal/config"
	"github.com/dusk-chancellor/dc-sso/pkg/zaplog"
)

func main() {
	cfg := config.MustLoad()

	log := zaplog.New()

	ctx := context.Background()

	app := app.New(ctx, log, cfg)

	go func() {
		app.GRPCServer.MustRun()
	}()

	// graceful shutdown

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop
	
	app.GRPCServer.Stop()
	log.Info("gracefully stopped server")
}
