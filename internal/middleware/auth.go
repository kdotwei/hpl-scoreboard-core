// package middleware

// import (
// 	"context"
// 	"net/http"
// 	"strings"

// 	"github.com/kdotwei/hpl-scoreboard/internal/token" // 引用剛剛建立的 package
// )

// // authorizationPayloadKey 必須是匯出(Exported)的，或者提供一個公開的常數
// // 建議直接定義一個公開的字串常數，簡單明瞭
// const AuthorizationPayloadKey = "authorization_payload"

// // 或者更嚴謹的 context key type (Week 11 進階寫法)
// // type contextKey string
// // const AuthorizationPayloadKey contextKey = "authorization_payload"

// func AuthMiddleware(tokenMaker token.Maker) func(http.Handler) http.Handler {
// 	return func(next http.Handler) http.Handler {
// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			authorizationHeader := r.Header.Get("Authorization")
// 			if len(authorizationHeader) == 0 {
// 				http.Error(w, "authorization header is not provided", http.StatusUnauthorized)
// 				return
// 			}

// 			fields := strings.Fields(authorizationHeader)
// 			if len(fields) < 2 {
// 				http.Error(w, "invalid authorization header format", http.StatusUnauthorized)
// 				return
// 			}

// 			authorizationType := strings.ToLower(fields[0])
// 			if authorizationType != "bearer" {
// 				http.Error(w, "unsupported authorization type", http.StatusUnauthorized)
// 				return
// 			}

// 			accessToken := fields[1]
// 			// 這裡假設你有實作 VerifyToken
// 			payload, err := tokenMaker.VerifyToken(accessToken)
// 			if err != nil {
// 				http.Error(w, "invalid token: "+err.Error(), http.StatusUnauthorized)
// 				return
// 			}

// 			// 將 Payload 塞入 Context
// 			ctx := context.WithValue(r.Context(), AuthorizationPayloadKey, payload)
// 			next.ServeHTTP(w, r.WithContext(ctx))
// 		})
// 	}
// }

package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/kdotwei/hpl-scoreboard/internal/token"
)

// 定義 Key
const AuthorizationPayloadKey = "authorization_payload"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. 檢查 Header
		if r.Header.Get("Authorization") == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// 2. 略過真實驗證，直接塞入一個有效的 Payload
		// 這樣 Handler 就能透過 ctx.Value(AuthorizationPayloadKey) 拿到資料
		mockPayload := &token.Payload{
			ID:        uuid.New(),
			Username:  "real-student-109704065", // 這裡換成你想測試的 ID
			IssuedAt:  time.Now(),
			ExpiredAt: time.Now().Add(time.Hour),
		}

		ctx := context.WithValue(r.Context(), AuthorizationPayloadKey, mockPayload)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
