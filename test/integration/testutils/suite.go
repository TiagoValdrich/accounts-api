package testutils

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/gofrs/uuid/v5"
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

const (
	databaseDoesNotExistsErrorCode = "3D000"
)

var migrationsFolder = "../../../db/migrations"

type testDBConfig struct {
	host       string
	user       string
	password   string
	dbName     string
	testDbName string
	sslMode    string
	debugMode  bool
}

type TestSuite struct {
	App *fiber.App
	DB  *bun.DB
}

func Setup() *TestSuite {
	dbCfg := createTestDatabaseIfNotExists()
	bunDB := connectToTestDatabase(dbCfg)
	db.RunMigrations(bunDB, &migrationsFolder)
	app := setupApp(bunDB)

	return &TestSuite{
		App: app,
		DB:  bunDB,
	}
}

func Teardown(testSuite *TestSuite) {
	if testSuite.DB != nil {
		if err := testSuite.DB.Close(); err != nil {
			panic(err)
		}

	}
}

func CleanupTables(t *testing.T, testSuite *TestSuite) {
	t.Helper()

	tables := []string{"transactions", "customer_account", "customer"}
	for _, table := range tables {
		_, err := testSuite.DB.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table))
		if err != nil {
			t.Fatalf("failed to truncate table %s: %v", table, err)
		}
	}
}

func createTestDatabaseIfNotExists() testDBConfig {
	cfg := getTestConfig()
	dbDsn := buildDatabaseConnectionString(cfg, false)

	adminDB := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dbDsn)))
	defer func() {
		if err := adminDB.Close(); err != nil {
			panic(err)
		}
	}()

	var exists bool
	err := adminDB.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)",
		cfg.testDbName,
	).Scan(&exists)

	if err != nil {
		var pgErr pgdriver.Error
		if errors.As(err, &pgErr); pgErr.Field('C') != databaseDoesNotExistsErrorCode {
			panic(fmt.Sprintf("unexpected error when checking database existence %v", err))
		}
	}

	if !exists {
		_, err = adminDB.Exec(fmt.Sprintf("CREATE DATABASE %s", cfg.testDbName))

		if err != nil {
			panic(fmt.Sprintf("failed to create test database: %v", err))
		}

		fmt.Printf("Created test database: %s\n", cfg.testDbName)
	}

	return cfg
}

func connectToTestDatabase(dbCfg testDBConfig) *bun.DB {
	dsn := buildDatabaseConnectionString(dbCfg, true)

	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	if err := sqldb.Ping(); err != nil {
		panic(fmt.Sprintf("failed to connect to test database: %v", err))
	}

	bunDB := bun.NewDB(sqldb, pgdialect.New())

	if dbCfg.debugMode {
		bunDB.AddQueryHook(bundebug.NewQueryHook(
			bundebug.WithVerbose(true),
			bundebug.WithEnabled(true),
		))
	}

	return bunDB
}

func buildDatabaseConnectionString(dbCfg testDBConfig, useTestDB bool) string {
	databaseName := dbCfg.dbName
	if useTestDB {
		databaseName = dbCfg.testDbName
	}

	return fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=%s",
		dbCfg.user, dbCfg.password, dbCfg.host, databaseName, dbCfg.sslMode,
	)
}

func getRandomizedDatabaseName() string {
	databaseID, err := uuid.NewV6()
	if err != nil {
		panic(fmt.Sprintf("failed to generate database db name ID: %v", err))
	}

	sanitizedID := strings.ReplaceAll(databaseID.String(), "-", "")

	return "accounts_api_test" + "_" + sanitizedID
}

func getTestConfig() testDBConfig {
	testDatabaseName := getRandomizedDatabaseName()

	return testDBConfig{
		host:       getEnvOrDefault("TEST_DB_HOST", "localhost:5432"),
		user:       getEnvOrDefault("TEST_DB_USER", "postgres"),
		password:   getEnvOrDefault("TEST_DB_PASSWORD", "postgres"),
		dbName:     getEnvOrDefault("TEST_DB_NAME", "postgres"),
		testDbName: testDatabaseName,
		sslMode:    getEnvOrDefault("TEST_DB_SSLMODE", "disable"),
		debugMode:  getEnvOrDefaultBool("TEST_DB_DEBUG", false),
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
	transactionRepository := repository.NewTransactionRepository(bunDB)

	accountsService := accounts.NewService(customerRepository, customerAccountRepository)
	accounts.NewHTTPHandler(router.GetApp(), accountsService)

	transactionsService := transactions.NewService(
		transactionRepository,
		customerAccountRepository,
	)
	transactions.NewHTTPHandler(router.GetApp(), transactionsService)

	return router.GetApp()
}
