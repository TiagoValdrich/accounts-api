package db

import (
	"fmt"

	"github.com/rs/zerolog/log"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/tiagovaldrich/accounts-api/internal/config"
	"github.com/uptrace/bun"
)

func RunMigrations(db *bun.DB) {
	migrations := &migrate.FileMigrationSource{
		Dir: config.MigrationsFolder,
	}

	n, err := migrate.Exec(db.DB, "postgres", migrations, migrate.Up)
	if err != nil {
		panic(fmt.Sprintf("failed to run migrations: %v", err))
	}

	if n == 0 {
		log.Info().Msg("no new migrations to run")
		return
	}

	log.Info().Msgf("applied %d migrations", n)
}
