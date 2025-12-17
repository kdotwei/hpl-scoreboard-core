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
		log.Fatalf("Unable to create connection pool: %v\n", err)
	}
	defer connPool.Close()

	// 關鍵：加入 Ping 來驗證實際連線
	err = connPool.Ping(context.Background())
	if err != nil {
		log.Fatalf("Unable to connect to database (Ping failed): %v\n", err)
	}

	log.Println("Connected to database successfully")

	// 3. 依賴注入 (Dependency Injection)
	store := db.New(connPool)
	svc := service.NewService(store)

	// 初始化 Token Maker
	// 注意：在正式環境中，這個 Secret Key 應該從環境變數讀取
	tokenMaker, err := token.NewJWTMaker("12345678901234567890123456789012")
	if err != nil {
		log.Fatal("cannot create token maker:", err)
	}

	// 注入 Service 和 TokenMaker
	h := handler.NewHandler(svc, tokenMaker)

	// 4. 路由設定 (Router)
	mux := http.NewServeMux()

	// [Route 1] Login (公開)
	mux.HandleFunc("POST /api/v1/login", h.Login)

	// [Route 2] Submit Score (需要 Auth)
	// 確保這一行只出現一次！
	// mux.Handle("POST /api/v1/scores", middleware.AuthMiddleware(http.HandlerFunc(h.CreateScore)))
	// 建立 Middleware 實例，注入 tokenMaker
	authMiddleware := middleware.AuthMiddleware(tokenMaker)
	mux.Handle("POST /api/v1/scores", authMiddleware(http.HandlerFunc(h.CreateScore)))

	// 5. 啟動伺服器
	log.Printf("Server starting on %s", serverAddress)
	if err := http.ListenAndServe(serverAddress, mux); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
