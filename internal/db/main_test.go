package db

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"

	// 👇 修正 Imports
	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var testStore Querier

func TestMain(m *testing.M) {
	ctx := context.Background()

	// 1. 啟動 Postgres Container
	pgContainer, err := tcpostgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:15-alpine"), // 👈 修正：使用 testcontainers.WithImage
		tcpostgres.WithDatabase("hpl_test"),
		tcpostgres.WithUsername("user"),
		tcpostgres.WithPassword("password"),
		testcontainers.WithWaitStrategy( // 👈 修正：使用 testcontainers.WithWaitStrategy
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		log.Fatalf("failed to start container: %s", err)
	}

	// 2. 取得連線字串
	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		log.Fatalf("failed to get connection string: %s", err)
	}

	// 3. 執行 Migration
	runDBMigration(connStr, "../../migrations")

	// 4. 連線 DB (使用 pgxpool)
	connPool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		log.Fatalf("failed to connect to db: %s", err)
	}
	defer connPool.Close()

	// 👇 這裡如果報錯 undefined New，是正常的 Red Phase (因為還沒 generate)
	// 但如果 generate 過了，加上 sqlc.yaml 的修正，這裡的型別錯誤就會消失
	testStore = New(connPool)

	code := m.Run()

	pgContainer.Terminate(ctx)
	os.Exit(code)
}

func runDBMigration(migrationURL string, sourceURL string) {
	db, err := sql.Open("postgres", migrationURL)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal("cannot create driver:", err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://"+sourceURL,
		"postgres", driver)
	if err != nil {
		log.Fatal("cannot create new migrate instance:", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("failed to run migrate up:", err)
	}
}
