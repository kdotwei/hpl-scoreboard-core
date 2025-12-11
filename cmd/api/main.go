package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kdotwei/hpl-scoreboard/internal/db"
	"github.com/kdotwei/hpl-scoreboard/internal/handler"
	"github.com/kdotwei/hpl-scoreboard/internal/middleware"
	"github.com/kdotwei/hpl-scoreboard/internal/service"
	"github.com/kdotwei/hpl-scoreboard/internal/token"
)

func main() {
	// 1. 環境變數設定
	// 建議之後整合到 internal/config
	dbSource := os.Getenv("DB_SOURCE")
	if dbSource == "" {
		dbSource = "postgresql://user:password@localhost:5432/hpl_scoreboard?sslmode=disable"
	}
	serverAddress := os.Getenv("SERVER_ADDRESS")
	if serverAddress == "" {
		serverAddress = ":8080"
	}

	// 2. 資料庫連線 (Database Layer)
	connPool, err := pgxpool.New(context.Background(), dbSource)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer connPool.Close()
	log.Println("Connected to database successfully")

	// 3. 依賴注入 (Dependency Injection)
	// Layer 3: Data Access
	store := db.New(connPool)
	// Layer 2: Business Logic
	svc := service.NewService(store)

	tokenMaker, err := token.NewJWTMaker("12345678901234567890123456789012")
	if err != nil {
		log.Fatal("cannot create token maker:", err)
	}

	// 注入 Service 和 TokenMaker
	h := handler.NewHandler(svc, tokenMaker)

	// 4. 路由設定 (Router) - 使用 Go 1.22+ 新語法
	mux := http.NewServeMux()

	// 註冊 Login 路由 (公開路由，不需要 AuthMiddleware)
	mux.HandleFunc("POST /api/v1/login", h.Login) // 👈 新增這行

	// 註冊 Score 路由 (保護路由)
	mux.Handle("POST /api/v1/scores", middleware.AuthMiddleware(http.HandlerFunc(h.CreateScore)))

	// 註冊路由 (Endpoint: POST /api/v1/scores)
	// 使用 Auth Middleware 保護此路由
	// 注意：這裡假設你的 middleware.AuthMiddleware 簽名是 func(http.Handler) http.Handler
	mux.Handle("POST /api/v1/scores", middleware.AuthMiddleware(http.HandlerFunc(h.CreateScore)))

	// 5. 啟動伺服器
	log.Printf("Server starting on %s", serverAddress)
	if err := http.ListenAndServe(serverAddress, mux); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
