package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kdotwei/hpl-scoreboard/internal/token"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware_MockBehavior(t *testing.T) {
	// 1. 定義一個用來捕捉 Context 結果的 "Next Handler"
	// 我們在這個 Handler 裡面檢查 Context 是否有被 Middleware 正確修改
	var capturedPayload *token.Payload

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 嘗試從 Context 取值
		val := r.Context().Value(AuthorizationPayloadKey)

		// 斷言：Context 裡必須要有東西
		assert.NotNil(t, val, "Context value should not be nil")

		// 斷言：型別必須正確
		payload, ok := val.(*token.Payload)
		assert.True(t, ok, "Value should be of type *token.Payload")

		// 捕捉起來供稍後檢查
		capturedPayload = payload

		w.WriteHeader(http.StatusOK)
	})

	// 2. 建立測試請求
	// 注意：雖然我們現在是用 Mock，但原本的邏輯還是會檢查 Header 是否存在
	req := httptest.NewRequest(http.MethodGet, "/api/v1/test", nil)
	req.Header.Set("Authorization", "Bearer any_token_here") // 必須有 Header 才能通過第一關

	rr := httptest.NewRecorder()

	// 3. 執行 Middleware
	// AuthMiddleware(nextHandler) 會回傳一個包裝後的 Handler
	AuthMiddleware(nextHandler).ServeHTTP(rr, req)

	// 4. 驗證結果
	assert.Equal(t, http.StatusOK, rr.Code)

	// 驗證 Payload 內容是否就是我們在 AuthMiddleware 裡寫死的那個 Mock 資料
	assert.NotNil(t, capturedPayload)
	assert.Equal(t, "real-student-109704065", capturedPayload.Username)
}
