package main

import (
	"github.com/tiagovaldrich/accounts-api/db"
	"github.com/tiagovaldrich/accounts-api/internal/config"
)

func main() {
	cfg := config.MustLoad()

	database := db.NewDatabase(&cfg.EnvVars.Database)
	defer database.Close()

	db.RunMigrations(database)
}
