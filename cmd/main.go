package main

import (
	"github.com/tiagovaldrich/accounts-api/db"
	"github.com/tiagovaldrich/accounts-api/internal/api/accounts"
	"github.com/tiagovaldrich/accounts-api/internal/config"
	"github.com/tiagovaldrich/accounts-api/internal/repository"
)

func main() {
	cfg := config.MustLoad()
	appRouter := config.NewRouter()

	database := db.NewDatabase(&cfg.EnvVars.Database)
	defer database.Close()

	db.RunMigrations(database)

	customerRepository := repository.NewCustomerRepository(database)
	customerAccountRepository := repository.NewCustomerAccountRepository(database)

	accountsService := accounts.NewService(customerRepository, customerAccountRepository)

	accounts.NewHTTPHandler(appRouter.GetApp(), accountsService)

	appRouter.Start()
}
