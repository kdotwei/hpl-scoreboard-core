package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateScore(t *testing.T) {
	// 由於我們還沒生成代碼，這裡的 CreateScoreParams 會報錯 (Red)
	arg := CreateScoreParams{
		UserID:       "user-uuid-mock", // 暫時使用假 ID
		Gflops:       123.45,
		ProblemSizeN: 10000,
		BlockSizeNb:  256,
		SubmittedAt:  time.Now(),
	}

	score, err := testStore.CreateScore(context.Background(), arg)

	assert.NoError(t, err)
	assert.NotEmpty(t, score)
	assert.Equal(t, arg.Gflops, score.Gflops)
	assert.NotZero(t, score.ID)
	assert.WithinDuration(t, arg.SubmittedAt, score.SubmittedAt, time.Second)
}
