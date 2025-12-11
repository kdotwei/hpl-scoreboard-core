package handler

import (
	"bytes"
	"context" // 👈 1. 新增
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time" // 👈 2. 新增 (為了初始化 token payload)

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/kdotwei/hpl-scoreboard/internal/db"
	"github.com/kdotwei/hpl-scoreboard/internal/middleware" // 👈 3. 新增
	"github.com/kdotwei/hpl-scoreboard/internal/service"
	"github.com/kdotwei/hpl-scoreboard/internal/service/mocks"
	"github.com/kdotwei/hpl-scoreboard/internal/token" // 👈 4. 新增
	token_mocks "github.com/kdotwei/hpl-scoreboard/internal/token/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateScore(t *testing.T) {
	// 1. Setup Mock
	mockService := new(mocks.Service)
	mockTokenMaker := new(token_mocks.Maker)     // 新增
	h := NewHandler(mockService, mockTokenMaker) // 修改這行

	reqBody := CreateScoreRequest{
		Gflops:       123.45,
		ProblemSizeN: 10000,
		BlockSizeNb:  256,
	}
	jsonBody, _ := json.Marshal(reqBody)

	// 定義 Mock User
	mockUser := "test-user"

	// 模擬 Service 行為
	// 注意：這裡可以順便驗證 UserID 是否正確傳遞
	mockService.On("CreateScore", mock.Anything, mock.MatchedBy(func(arg service.CreateScoreParams) bool {
		return arg.Gflops == 123.45 && arg.UserID == mockUser // 👈 加上 UserID 驗證
	})).Return(&db.Score{
		ID:     pgtype.UUID{Bytes: [16]byte{}, Valid: true},
		Gflops: 123.45,
	}, nil)

	// 2. 建立 HTTP Request
	req, _ := http.NewRequest("POST", "/api/v1/scores", bytes.NewBuffer(jsonBody))

	// ✨✨✨ 關鍵修正：注入 Auth Payload 到 Context ✨✨✨
	mockPayload := &token.Payload{
		Username:  mockUser,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(time.Minute),
	}
	// 模擬 Middleware 的行為，將 Payload 塞入 Context
	ctx := context.WithValue(req.Context(), middleware.AuthorizationPayloadKey, mockPayload)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	// 3. 執行 Handler
	http.HandlerFunc(h.CreateScore).ServeHTTP(rr, req)

	// 4. Assertions
	assert.Equal(t, http.StatusCreated, rr.Code)

	mockService.AssertExpectations(t)
}
