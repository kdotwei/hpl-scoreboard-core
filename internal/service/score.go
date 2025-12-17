package service

import (
	"context"
	"time"

	"github.com/kdotwei/hpl-scoreboard/internal/db"
)

func (s *HPLService) CreateScore(ctx context.Context, arg CreateScoreParams) (*db.Score, error) {
	result, err := s.store.CreateScore(ctx, db.CreateScoreParams{
		UserID:        arg.UserID,
		Gflops:        arg.Gflops,
		ProblemSizeN:  int32(arg.ProblemSizeN),
		BlockSizeNb:   int32(arg.BlockSizeNb),
		LinuxUsername: arg.LinuxUsername,
		N:             int32(arg.N),
		Nb:            int32(arg.NB), // 修正編譯錯誤：sqlc 生成的是 Nb
		P:             int32(arg.P),
		Q:             int32(arg.Q),
		ExecutionTime: arg.ExecutionTime,
		SubmittedAt:   time.Now(), // 確保帶上時間戳記
	})
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *HPLService) ListScores(ctx context.Context, limit int32, offset int32) ([]db.Score, error) {
	return s.store.ListTopScores(ctx, db.ListTopScoresParams{
		Limit:  limit,
		Offset: offset,
	})
}

func (s *HPLService) ListScoresWithPagination(ctx context.Context, params ListScoresParams) (*PaginatedScoresResponse, error) {
	// Get scores with pagination
	scores, err := s.store.ListTopScores(ctx, db.ListTopScoresParams{
		Limit:  params.Limit,
		Offset: params.Offset,
	})
	if err != nil {
		return nil, err
	}

	// Get total count for frontend reference
	totalRecords, err := s.store.CountTotalScores(ctx)
	if err != nil {
		return nil, err
	}

	// Calculate if there are more records
	hasMore := int64(params.Offset+int32(len(scores))) < totalRecords

	response := &PaginatedScoresResponse{
		Scores:       scores,
		HasMore:      hasMore,
		TotalRecords: totalRecords,
		Limit:        params.Limit,
		Offset:       params.Offset,
	}

	return response, nil
}
