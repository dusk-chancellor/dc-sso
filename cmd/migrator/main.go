package main

import (
	"context"
	"database/sql"
	"flag"
	"log"

	"github.com/dusk-chancellor/dc-sso/internal/config"
	"github.com/dusk-chancellor/dc-sso/internal/database/postgres"
	m "github.com/dusk-chancellor/dc-sso/migrations"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

// migrator runs migrations independently from sso app
// so it would be more managable
// now only 'up' & 'down' commands supported

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	pool, err := postgres.ConnectDB(context.Background(), &cfg.Db)
	if err != nil {
		panic(err)
	}

	goose.SetDialect("postgres")
	goose.SetBaseFS(m.Migrations.FS)

	db := stdlib.OpenDBFromPool(pool)

	migrate := fetchMigrate()
	switch migrate {
	case "up":
		err = up(db, m.Migrations.Dir)
	case "down":
		err = down(db, m.Migrations.Dir)
	default:
		panic("unknown command")
	}
	if err != nil {
		panic(err)
	}

	log.Print("Successfully ran migrations!\n")
}
// getting flag value for executing command
// e.g -> `migrator --migrate="up"`
func fetchMigrate() string {
	var res string

	flag.StringVar(&res, "migrate", "", "migration command")
	flag.Parse()

	return res
}
// `goose up`
func up(db *sql.DB, dir string) error {
	return goose.Up(db, dir)
}
// `goose down`
func down(db *sql.DB, dir string) error {
	return goose.Down(db, dir)
}
