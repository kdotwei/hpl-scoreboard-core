package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kdotwei/hpl-scoreboard/internal/service/mocks"
	token_mocks "github.com/kdotwei/hpl-scoreboard/internal/token/mocks" // 引用剛剛生成的 token mock
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLogin(t *testing.T) {
	// 1. Setup Mocks
	mockService := new(mocks.Service)
	mockTokenMaker := new(token_mocks.Maker) // 新增 TokenMaker Mock

	// 🔴 這裡會報錯：因為目前的 NewHandler 只接受 service，不接受 tokenMaker
	h := NewHandler(mockService, mockTokenMaker)

	// 2. 準備 Request
	user := "agent-lead"
	reqBody := LoginRequest{ // 🔴 這裡會報錯：LoginRequest 尚未定義
		Username: user,
	}
	jsonBody, _ := json.Marshal(reqBody)

	// 3. 設定 Mock 行為
	// 當呼叫 CreateToken 時，回傳一個假 Token
	mockTokenMaker.On("CreateToken", user, mock.Anything).Return("mock_access_token", nil, nil)

	req, _ := http.NewRequest("POST", "/api/v1/login", bytes.NewBuffer(jsonBody))
	rr := httptest.NewRecorder()

	// 4. 執行 Handler
	// 🔴 這裡會報錯：Login 方法尚未實作
	http.HandlerFunc(h.Login).ServeHTTP(rr, req)

	// 5. 驗證
	assert.Equal(t, http.StatusOK, rr.Code)

	// 驗證回傳的 JSON 包含 access_token
	var resp LoginResponse // 🔴 這裡會報錯：LoginResponse 尚未定義
	// Fix errcheck: 檢查 Decode 錯誤
	err := json.NewDecoder(rr.Body).Decode(&resp)
	assert.NoError(t, err) // 加上這行斷言
	assert.Equal(t, "mock_access_token", resp.AccessToken)
	assert.Equal(t, user, resp.User.Username) // 假設我們也會回傳 User 資訊
}
