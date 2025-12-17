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
		if parsedLimit, err := strconv.ParseInt(limitStr, 10, 32); err != nil || parsedLimit <= 0 {
			http.Error(w, "Invalid limit parameter", http.StatusBadRequest)
			return
		} else {
			limit = int32(parsedLimit)
		}
	}

	// Parse optional offset query parameter (default to 0)
	offset := int32(0)
	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if parsedOffset, err := strconv.ParseInt(offsetStr, 10, 32); err != nil || parsedOffset < 0 {
			http.Error(w, "Invalid offset parameter", http.StatusBadRequest)
			return
		} else {
			offset = int32(parsedOffset)
		}
	}

	// Get scores from service
	scores, err := h.service.ListScores(r.Context(), limit, offset)
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

func (h *Handler) ListScoresWithPagination(w http.ResponseWriter, r *http.Request) {
	// Parse pagination parameters
	params := service.ListScoresParams{
		Limit:  10, // Default limit
		Offset: 0,  // Default offset
	}

	// Parse limit query parameter
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if parsedLimit, err := strconv.ParseInt(limitStr, 10, 32); err != nil || parsedLimit <= 0 || parsedLimit > 100 {
			http.Error(w, "Invalid limit parameter (must be 1-100)", http.StatusBadRequest)
			return
		} else {
			params.Limit = int32(parsedLimit)
		}
	}

	// Parse offset query parameter
	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if parsedOffset, err := strconv.ParseInt(offsetStr, 10, 32); err != nil || parsedOffset < 0 {
			http.Error(w, "Invalid offset parameter (must be >= 0)", http.StatusBadRequest)
			return
		} else {
			params.Offset = int32(parsedOffset)
		}
	}

	// Get paginated scores from service
	response, err := h.service.ListScoresWithPagination(r.Context(), params)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
