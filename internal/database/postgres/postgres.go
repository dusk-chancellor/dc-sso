package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/dusk-chancellor/dc-sso/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

// db connection pool setup

// creates pool connection w/ db
func ConnectDB(ctx context.Context, cfg *config.DB) (*pgxpool.Pool, error) {
	dsn := buildDSN(cfg) // db url

	poolCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	// minimalistic settings for pool
	poolCfg.MaxConns = 10 // pool max size
	poolCfg.MinConns = 2 // pool min size
	poolCfg.MaxConnIdleTime = 15 * time.Minute
	poolCfg.MaxConnLifetime = 30 * time.Minute

	// creating new pool connection
	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, err
	}

	// pinging db pool
	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	return pool, nil
}
// db url construction
func buildDSN(cfg *config.DB) string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=%s", 
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
		cfg.SSLMode,
	)
}
