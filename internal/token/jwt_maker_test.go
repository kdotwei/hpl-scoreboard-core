package token

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
)

func TestJWTMaker(t *testing.T) {
	// 定義一個足夠長的 Secret Key (至少 32 bytes)
	secretKey := "12345678901234567890123456789012"
	maker, err := NewJWTMaker(secretKey)
	require.NoError(t, err)

	username := "test-user"
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	// 1. 測試建立 Token
	token, payload, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	// 2. 測試驗證 Token
	payload, err = maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

func TestExpiredJWTToken(t *testing.T) {
	secretKey := "12345678901234567890123456789012"
	maker, err := NewJWTMaker(secretKey)
	require.NoError(t, err)

	// 建立一個「負時間」的 Token (立刻過期)
	token, payload, err := maker.CreateToken("test-user", -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	// 驗證應該要失敗
	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	// 注意：jwt/v5 回傳的錯誤可能包含 "token is expired" 字串
	require.Nil(t, payload)
}

func TestInvalidJWTTokenAlgNone(t *testing.T) {
	// 模擬一個駭客攻擊：使用 "None" 演算法簽署 Token
	// 這是一種常見的 JWT 漏洞，我們的 VerifyToken 必須要能擋下來
	payload, err := NewPayload("hacker", time.Minute)
	require.NoError(t, err)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	// 使用 UnsafeAllowNoneSignatureType 來模擬未簽名的 Token
	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	secretKey := "12345678901234567890123456789012"
	maker, err := NewJWTMaker(secretKey)
	require.NoError(t, err)

	// 驗證應該要失敗，並回傳 ErrInvalidToken
	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrInvalidToken.Error())
	require.Nil(t, payload)
}
