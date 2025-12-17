package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/kdotwei/hpl-scoreboard/internal/token"
)

type contextKey string

const AuthorizationPayloadKey contextKey = "authorization_payload"

// AuthMiddleware 改為回傳一個 Closure，因為它需要依賴 tokenMaker
func AuthMiddleware(tokenMaker token.Maker) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 1. 取得 Header
			authorizationHeader := r.Header.Get("Authorization")
			if len(authorizationHeader) == 0 {
				http.Error(w, "authorization header is not provided", http.StatusUnauthorized)
				return
			}

			// 2. 解析 Bearer Token 格式
			fields := strings.Fields(authorizationHeader)
			if len(fields) < 2 {
				http.Error(w, "invalid authorization header format", http.StatusUnauthorized)
				return
			}

			authorizationType := strings.ToLower(fields[0])
			if authorizationType != "bearer" {
				http.Error(w, "unsupported authorization type", http.StatusUnauthorized)
				return
			}

			accessToken := fields[1]

			// 3. ✨ 關鍵改變：使用 Maker 進行真實驗證 ✨
			payload, err := tokenMaker.VerifyToken(accessToken)
			if err != nil {
				http.Error(w, "invalid token: "+err.Error(), http.StatusUnauthorized)
				return
			}

			// 4. 將解析出來的 Payload (包含 username) 塞入 Context
			ctx := context.WithValue(r.Context(), AuthorizationPayloadKey, payload)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
