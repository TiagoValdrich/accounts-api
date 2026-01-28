package integration

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/tiagovaldrich/accounts-api/db"
	"github.com/tiagovaldrich/accounts-api/internal/api/accounts"
	"github.com/tiagovaldrich/accounts-api/internal/api/transactions"
	"github.com/tiagovaldrich/accounts-api/internal/config"
	"github.com/tiagovaldrich/accounts-api/internal/repository"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

type testDBConfig struct {
	host      string
	user      string
	password  string
	dbName    string
	sslMode   string
	debugMode bool
}

var (
	App              *fiber.App
	DB               *bun.DB
	migrationsFolder string = "../../db/migrations"
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func setup() {
	createTestDatabaseIfNotExists()
	DB = connectToTestDatabase()
	db.RunMigrations(DB, &migrationsFolder)
	App = setupApp(DB)
}

func teardown() {
	if DB != nil {
		DB.Close()
	}
}

func createTestDatabaseIfNotExists() {
	cfg := getTestConfig()

	adminDSN := fmt.Sprintf(
		"postgres://%s:%s@%s/postgres?sslmode=%s",
		cfg.user, cfg.password, cfg.host, cfg.sslMode,
	)

	adminDB := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(adminDSN)))
	defer adminDB.Close()

	var exists bool
	err := adminDB.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)",
		cfg.dbName,
	).Scan(&exists)
	if err != nil {
		panic(fmt.Sprintf("failed to check if test database exists: %v", err))
	}

	if !exists {
		_, err = adminDB.Exec(fmt.Sprintf("CREATE DATABASE %s", cfg.dbName))
		if err != nil {
			panic(fmt.Sprintf("failed to create test database: %v", err))
		}
		fmt.Printf("Created test database: %s\n", cfg.dbName)
	}
}

func connectToTestDatabase() *bun.DB {
	cfg := getTestConfig()
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=%s",
		cfg.user, cfg.password, cfg.host, cfg.dbName, cfg.sslMode,
	)

	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	if err := sqldb.Ping(); err != nil {
		panic(fmt.Sprintf("failed to connect to test database: %v", err))
	}

	bunDB := bun.NewDB(sqldb, pgdialect.New())

	if cfg.debugMode {
		bunDB.AddQueryHook(bundebug.NewQueryHook(
			bundebug.WithVerbose(true),
			bundebug.WithEnabled(true),
		))
	}

	return bunDB
}

func getTestConfig() testDBConfig {
	return testDBConfig{
		host:      getEnvOrDefault("TEST_DB_HOST", "localhost:5432"),
		user:      getEnvOrDefault("TEST_DB_USER", "postgres"),
		password:  getEnvOrDefault("TEST_DB_PASSWORD", "postgres"),
		dbName:    getEnvOrDefault("TEST_DB_NAME", "accounts_api_test"),
		sslMode:   getEnvOrDefault("TEST_DB_SSLMODE", "disable"),
		debugMode: getEnvOrDefaultBool("TEST_DB_DEBUG", false),
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvOrDefaultBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		return value == "true" || value == "1"
	}
	return defaultValue
}

func setupApp(bunDB *bun.DB) *fiber.App {
	router := config.NewRouter()

	customerRepository := repository.NewCustomerRepository(bunDB)
	customerAccountRepository := repository.NewCustomerAccountRepository(bunDB)
	balanceRepository := repository.NewBalanceRepository(bunDB)
	transactionRepository := repository.NewTransactionRepository(bunDB)

	accountsService := accounts.NewService(customerRepository, customerAccountRepository, balanceRepository)
	accounts.NewHTTPHandler(router.GetApp(), accountsService)

	transactionsService := transactions.NewService(
		transactionRepository,
		customerAccountRepository,
		balanceRepository,
	)
	transactions.NewHTTPHandler(router.GetApp(), transactionsService)

	return router.GetApp()
}

func CleanupTables(t *testing.T) {
	t.Helper()

	tables := []string{"transactions", "balance", "customer_account", "customer"}
	for _, table := range tables {
		_, err := DB.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table))
		if err != nil {
			t.Fatalf("failed to truncate table %s: %v", table, err)
		}
	}
}
