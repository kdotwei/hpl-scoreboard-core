package handler

import (
	"encoding/json"
	"net/http"

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

	authPayload := r.Context().Value(middleware.AuthorizationPayloadKey).(*token.Payload)

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
	json.NewEncoder(w).Encode(score)
}
