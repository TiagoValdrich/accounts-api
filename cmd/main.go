package main

import (
	"github.com/tiagovaldrich/accounts-api/db"
	"github.com/tiagovaldrich/accounts-api/internal/api/accounts"
	"github.com/tiagovaldrich/accounts-api/internal/api/transactions"
	"github.com/tiagovaldrich/accounts-api/internal/config"
	"github.com/tiagovaldrich/accounts-api/internal/repository"
)

func main() {
	cfg := config.MustLoad()
	appRouter := config.NewRouter()

	database := db.NewDatabase(&cfg.EnvVars.Database)
	defer database.Close()

	db.RunMigrations(database, nil)

	customerRepository := repository.NewCustomerRepository(database)
	customerAccountRepository := repository.NewCustomerAccountRepository(database)
	balanceRepository := repository.NewBalanceRepository(database)
	transactionRepository := repository.NewTransactionRepository(database)

	accountsService := accounts.NewService(customerRepository, customerAccountRepository, balanceRepository)
	transactionsService := transactions.NewService(transactionRepository, customerAccountRepository, balanceRepository)

	accounts.NewHTTPHandler(appRouter.GetApp(), accountsService)
	transactions.NewHTTPHandler(appRouter.GetApp(), transactionsService)

	appRouter.Start()
}
