package postgres

import (
	"database/sql"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func NewPostgresDB(cfg *Config) (*bun.DB, error) {
	sqldb := sql.OpenDB(pgdriver.NewConnector(
		pgdriver.WithNetwork("tcp"),
		pgdriver.WithAddr(cfg.Addr),
		pgdriver.WithUser(cfg.User),
		pgdriver.WithPassword(cfg.Pass),
		pgdriver.WithDatabase(cfg.Name),
		pgdriver.WithTLSConfig(nil),
	))

	db := bun.NewDB(sqldb, pgdialect.New())

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
