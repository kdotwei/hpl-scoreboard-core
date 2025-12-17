package service

import (
	"context"

	"github.com/kdotwei/hpl-scoreboard/internal/db"
)

// CreateScoreParams 是 Service 層的輸入參數
type CreateScoreParams struct {
	UserID        string
	Gflops        float64
	ProblemSizeN  int
	BlockSizeNb   int
	LinuxUsername string
	N             int
	NB            int // Service 層可以維持使用 NB 命名
	P             int
	Q             int
	ExecutionTime float64
}

// Service 定義了業務邏輯的介面
type Service interface {
	CreateScore(ctx context.Context, arg CreateScoreParams) (*db.Score, error)
	ListScores(ctx context.Context, limit int32) ([]db.Score, error)
}

// Ensure implementation (編譯時期檢查，確保 HPLService 有實作 Service)
// var _ Service = (*HPLService)(nil)

type HPLService struct {
	store db.Querier
}

func NewService(store db.Querier) *HPLService {
	return &HPLService{store: store}
}
