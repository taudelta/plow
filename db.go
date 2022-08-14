package plow

import (
	"database/sql"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/taudelta/plow/fixtures"
)

type DbConfig struct {
	DSN string
}

type storages struct {
	Postgres *sql.DB
}

var storageRegistry storages

func UsePostgresDB(cfg *DbConfig) {
	connConfig, err := pgx.ParseConfig(cfg.DSN)
	if err != nil {
		log.Fatal(err)
	}

	db := stdlib.OpenDB(*connConfig)

	storageRegistry.Postgres = db
}

func LoadPostgresFixtures(location string, names []string) {
	loader := fixtures.Postgres(storageRegistry.Postgres, location, false)
	if err := loader.Load(names); err != nil {
		log.Fatal(err)
	}
}
