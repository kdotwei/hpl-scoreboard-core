package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/kdotwei/hpl-scoreboard/internal/db"
	"github.com/kdotwei/hpl-scoreboard/internal/handler"
	"github.com/kdotwei/hpl-scoreboard/internal/middleware"
	"github.com/kdotwei/hpl-scoreboard/internal/service"
	"github.com/kdotwei/hpl-scoreboard/internal/token"
)

// enableCORS middleware allows cross-origin requests from frontend
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow specific origins
		origin := r.Header.Get("Origin")
		if origin == "http://localhost:5173" || origin == "http://localhost:3000" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight OPTIONS request
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	// 載入 .env 檔案
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using environment variables or defaults")
	}

	// 1. 環境變數設定
	dbSource := os.Getenv("DB_SOURCE")
	if dbSource == "" {
		dbSource = "postgresql://user:password@localhost:5432/hpl_scoreboard?sslmode=disable"
	}
	serverAddress := os.Getenv("SERVER_ADDRESS")
	if serverAddress == "" {
		serverAddress = ":8080"
	}

	// JWT Secret Key
	jwtSecretKey := os.Getenv("JWT_SECRET_KEY")
	if jwtSecretKey == "" {
		jwtSecretKey = "12345678901234567890123456789012" // 預設值（僅用於開發）
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
	tokenMaker, err := token.NewJWTMaker(jwtSecretKey)
	if err != nil {
		log.Fatal("cannot create token maker:", err)
	}

	// 注入 Service 和 TokenMaker
	h := handler.NewHandler(svc, tokenMaker)

	// 4. 路由設定 (Router)
	mux := http.NewServeMux()

	// [Route 1] Login (公開)
	mux.HandleFunc("POST /api/v1/login", h.Login)

	// [Route 2] List Scores (公開)
	mux.HandleFunc("GET /api/v1/scores", h.ListScores)

	// [Route 2.1] List Scores with Pagination (公開)
	mux.HandleFunc("GET /api/v1/scores/paginated", h.ListScoresWithPagination)

	// [Route 3] Submit Score (需要 Auth)
	authMiddleware := middleware.AuthMiddleware(tokenMaker)
	mux.Handle("POST /api/v1/scores", authMiddleware(http.HandlerFunc(h.CreateScore)))

	// 5. 啟動伺服器
	log.Printf("Server starting on %s", serverAddress)
	if err := http.ListenAndServe(serverAddress, enableCORS(mux)); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
