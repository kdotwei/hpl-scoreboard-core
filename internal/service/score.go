package service

import (
	"context"

	"github.com/kdotwei/hpl-scoreboard/internal/db"
)

func (s *HPLService) CreateScore(ctx context.Context, arg CreateScoreParams) (*db.Score, error) {
	// 1. 先取得數值結果
	result, err := s.store.CreateScore(ctx, db.CreateScoreParams{
		UserID:       arg.UserID,
		Gflops:       arg.Gflops,
		ProblemSizeN: int32(arg.ProblemSizeN),
		BlockSizeNb:  int32(arg.BlockSizeNb),
	})

	if err != nil {
		return nil, err
	}

	// 2. 回傳該數值的記憶體地址 (&result)
	return &result, nil
}
