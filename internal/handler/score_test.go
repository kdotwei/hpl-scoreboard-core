package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kdotwei/hpl-scoreboard/internal/db"
	"github.com/kdotwei/hpl-scoreboard/internal/service"
	"github.com/kdotwei/hpl-scoreboard/internal/service/mocks" // 👈 剛剛生成的
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateScore(t *testing.T) {
	// 1. Setup Mock
	mockService := new(mocks.Service)
	h := NewHandler(mockService) // 🔴 這裡會報錯，因為 NewHandler 還沒寫

	// 準備 Request Body
	reqBody := CreateScoreRequest{ // 🔴 這裡會報錯，因為 Struct 還沒定義
		Gflops:       123.45,
		ProblemSizeN: 10000,
		BlockSizeNb:  256,
	}
	jsonBody, _ := json.Marshal(reqBody)

	// 模擬 Service 行為：預期會被呼叫一次，並回傳成功
	mockService.On("CreateScore", mock.Anything, mock.MatchedBy(func(arg service.CreateScoreParams) bool {
		return arg.Gflops == 123.45 // 驗證參數傳遞正確
	})).Return(&db.Score{
		ID:     [16]byte{},
		Gflops: 123.45,
	}, nil)

	// 2. 建立 HTTP Request
	req, _ := http.NewRequest("POST", "/api/v1/scores", bytes.NewBuffer(jsonBody))
	rr := httptest.NewRecorder()

	// 3. 執行 Handler
	http.HandlerFunc(h.CreateScore).ServeHTTP(rr, req) // 🔴 CreateScore 方法還沒寫

	// 4. Assertions (斷言)
	assert.Equal(t, http.StatusCreated, rr.Code)

	mockService.AssertExpectations(t)
}
