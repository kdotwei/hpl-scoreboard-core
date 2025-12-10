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
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var testStore Store // sqlc 將會生成這個 Store 介面

func TestMain(m *testing.M) {
	ctx := context.Background()

	// 1. 啟動 Postgres Container
	pgContainer, err := tcpostgres.RunContainer(ctx,
		tcpostgres.WithImage("postgres:15-alpine"),
		tcpostgres.WithDatabase("hpl_test"),
		tcpostgres.WithUsername("user"),
		tcpostgres.WithPassword("password"),
		tcpostgres.WithWaitStrategy(
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
	// 修正路徑：從 internal/db 往上兩層找到 migrations 資料夾
	runDBMigration(connStr, "../../migrations")

	// 4. 連線 DB
	connPool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		log.Fatalf("failed to connect to db: %s", err)
	}
	defer connPool.Close()

	testStore = New(connPool) // sqlc 將會生成 New 函式

	// 5. 執行測試
	code := m.Run()

	// 6. 清理
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
