package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/kdotwei/hpl-scoreboard/internal/token"
)

// 👇 1. 定義自訂型別 (解決 SA1029)
type contextKey string

// 👇 2. 使用該型別定義常數
const AuthorizationPayloadKey contextKey = "authorization_payload"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. 檢查 Header
		if r.Header.Get("Authorization") == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// 2. 略過真實驗證，直接塞入一個有效的 Payload
		mockPayload := &token.Payload{
			ID:        uuid.New(),
			Username:  "real-student-109704065",
			IssuedAt:  time.Now(),
			ExpiredAt: time.Now().Add(time.Hour),
		}

		ctx := context.WithValue(r.Context(), AuthorizationPayloadKey, mockPayload)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
