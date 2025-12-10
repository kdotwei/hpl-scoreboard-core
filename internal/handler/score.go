package handler

import (
	"encoding/json"
	"net/http"

	"github.com/kdotwei/hpl-scoreboard/internal/middleware"
	"github.com/kdotwei/hpl-scoreboard/internal/service"
	"github.com/kdotwei/hpl-scoreboard/internal/token"
)

// CreateScoreRequest 定義前端/Agent 傳來的 JSON 格式
type CreateScoreRequest struct {
	Gflops       float64 `json:"gflops"`
	ProblemSizeN int     `json:"problem_size_n"`
	BlockSizeNb  int     `json:"block_size_nb"`
}

func (h *Handler) CreateScore(w http.ResponseWriter, r *http.Request) {
	// 1. 解析 Request Body
	var req CreateScoreRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// 2. 取得 UserID
	authPayload, ok := r.Context().Value(middleware.AuthorizationPayloadKey).(*token.Payload)
	if !ok {
		// 如果 Context 裡拿不到 Payload，或者是型別不對，回傳 401
		http.Error(w, "Unauthorized: missing token info", http.StatusUnauthorized)
		return
	}

	userID := authPayload.Username // 或是 authPayload.UserID，看你的 Token 結構定義

	// 3. 呼叫 Service 層處理業務邏輯
	score, err := h.service.CreateScore(r.Context(), service.CreateScoreParams{
		UserID:       userID,
		Gflops:       req.Gflops,
		ProblemSizeN: req.ProblemSizeN,
		BlockSizeNb:  req.BlockSizeNb,
	})

	if err != nil {
		// 實際專案建議使用 structured logging 紀錄錯誤
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// 4. 回傳成功 (201 Created) 與 JSON 結果
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(score)
}
