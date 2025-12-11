package handler

import (
	"encoding/json"
	"net/http"
	"time"
)

// LoginRequest 定義請求格式
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
}

// LoginResponse 定義回傳格式
type LoginResponse struct {
	AccessToken string       `json:"access_token"`
	User        UserResponse `json:"user"`
}

type UserResponse struct {
	Username string `json:"username"`
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// 由於沒有 User DB，我們直接簽發 Token
	// 設定 Token 有效期為 24 小時
	accessToken, _, err := h.tokenMaker.CreateToken(req.Username, 24*time.Hour)
	if err != nil {
		http.Error(w, "Failed to create access token", http.StatusInternalServerError)
		return
	}

	// 回傳結果
	resp := LoginResponse{
		AccessToken: accessToken,
		// Fix S1016: 直接將 req 轉型為 UserResponse
		// 由於兩個 Struct 欄位名稱與型別一致 (Username string)，Go 允許直接轉型 (忽略 Tags)
		User: UserResponse(req),
	}

	w.WriteHeader(http.StatusOK)

	// Fix errcheck: 檢查 Encode 錯誤 (保留您之前修復的邏輯)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		return
	}
}
