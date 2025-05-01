package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenMaker interface {
	CreateToken(userID uint, nickname string, email string, token_type string, expiresIn time.Duration) (string, error)
	VerifyToken(token string) (*Payload, error)
}

type Payload struct {
	jwt.RegisteredClaims
	UserID string `json:"user_id"`
	Nickname string `json:"nickname"`
	Email string `json:"email"`
	TokenType string `json:"token_type"`
	IssuedAt int64 `json:"issued_at"`
	ExpiresAt int64 `json:"expires_at"`
}

func NewPayload(userID, nickname string, email string, TokenType string, IssuedAt int64,ExpiresAt int64) *Payload {
	now := time.Now()
	return &Payload{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(24 * time.Hour)),
			IssuedAt: jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer: "anime-service",
		},
		UserID: userID,
		Nickname: nickname,
		Email: email,
	}
}