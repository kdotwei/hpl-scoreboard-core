package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jackc/pgx/v5/pgtype" // 👈 1. 新增這個 import
	"github.com/kdotwei/hpl-scoreboard/internal/db"
	"github.com/kdotwei/hpl-scoreboard/internal/service"
	"github.com/kdotwei/hpl-scoreboard/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateScore(t *testing.T) {
	// 1. Setup Mock
	mockService := new(mocks.Service)
	h := NewHandler(mockService)

	reqBody := CreateScoreRequest{
		Gflops:       123.45,
		ProblemSizeN: 10000,
		BlockSizeNb:  256,
	}
	jsonBody, _ := json.Marshal(reqBody)

	// 模擬 Service 行為
	mockService.On("CreateScore", mock.Anything, mock.MatchedBy(func(arg service.CreateScoreParams) bool {
		return arg.Gflops == 123.45
	})).Return(&db.Score{
		// 👇 2. 修正這裡：使用 pgtype.UUID
		ID:     pgtype.UUID{Bytes: [16]byte{}, Valid: true},
		Gflops: 123.45,
	}, nil)

	// 2. 建立 HTTP Request
	req, _ := http.NewRequest("POST", "/api/v1/scores", bytes.NewBuffer(jsonBody))
	rr := httptest.NewRecorder()

	// 3. 執行 Handler
	http.HandlerFunc(h.CreateScore).ServeHTTP(rr, req)

	// 4. Assertions
	assert.Equal(t, http.StatusCreated, rr.Code)

	mockService.AssertExpectations(t)
}
