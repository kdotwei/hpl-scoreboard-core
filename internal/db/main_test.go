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

	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var testStore Querier

func TestMain(m *testing.M) {
	ctx := context.Background()

	// 1. 啟動 Postgres Container
	// Fix SA1019: 使用 tcpostgres.Run 取代 RunContainer
	// 注意：Run 的第二個參數是 image name，所以移除了 testcontainers.WithImage
	pgContainer, err := tcpostgres.Run(ctx,
		"postgres:15-alpine",
		tcpostgres.WithDatabase("hpl_test"),
		tcpostgres.WithUsername("user"),
		tcpostgres.WithPassword("password"),
		testcontainers.WithWaitStrategy(
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

	testStore = New(connPool)

	code := m.Run()

	// 6. 清理
	// Fix errcheck: 處理 Terminate 的回傳錯誤
	if err := pgContainer.Terminate(ctx); err != nil {
		log.Printf("failed to terminate container: %s", err)
	}

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
