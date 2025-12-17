package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/kdotwei/hpl-scoreboard/internal/middleware"
	"github.com/kdotwei/hpl-scoreboard/internal/service"
	"github.com/kdotwei/hpl-scoreboard/internal/token"
)

type CreateScoreRequest struct {
	Gflops        float64 `json:"gflops"`
	ProblemSizeN  int     `json:"problem_size_n"`
	BlockSizeNb   int     `json:"block_size_nb"`
	LinuxUsername string  `json:"linux_username"`
	N             int     `json:"n"`
	NB            int     `json:"nb"`
	P             int     `json:"p"`
	Q             int     `json:"q"`
	ExecutionTime float64 `json:"execution_time"`
}

func (h *Handler) CreateScore(w http.ResponseWriter, r *http.Request) {
	var req CreateScoreRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	authPayloadValue := r.Context().Value(middleware.AuthorizationPayloadKey)
	if authPayloadValue == nil {
		http.Error(w, "Missing authorization payload", http.StatusUnauthorized)
		return
	}

	authPayload, ok := authPayloadValue.(*token.Payload)
	if !ok {
		http.Error(w, "Invalid authorization payload", http.StatusUnauthorized)
		return
	}

	score, err := h.service.CreateScore(r.Context(), service.CreateScoreParams{
		UserID:        authPayload.Username,
		Gflops:        req.Gflops,
		ProblemSizeN:  req.ProblemSizeN,
		BlockSizeNb:   req.BlockSizeNb,
		LinuxUsername: req.LinuxUsername,
		N:             req.N,
		NB:            req.NB,
		P:             req.P,
		Q:             req.Q,
		ExecutionTime: req.ExecutionTime,
	})

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(score); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) ListScores(w http.ResponseWriter, r *http.Request) {
	// Parse optional limit query parameter (default to 10)
	limit := int32(10)
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if parsedLimit, err := strconv.ParseInt(limitStr, 10, 32); err == nil && parsedLimit > 0 {
			limit = int32(parsedLimit)
		}
	}

	// Get scores from service
	scores, err := h.service.ListScores(r.Context(), limit)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(scores); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
