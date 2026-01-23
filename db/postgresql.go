package db

import (
	"database/sql"

	"github.com/tiagovaldrich/accounts-api/internal/config"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func NewDatabase(cfg *config.DatabaseConfig) *bun.DB {
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(cfg.DSN())))
	db := bun.NewDB(sqldb, pgdialect.New())

	return db
}
