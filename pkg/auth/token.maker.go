package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenMaker interface {
	CreateToken(userID uint, nickname string, email string) (string, error)
	VerifyToken(token string) (*Payload, error)
}

type Payload struct {
	jwt.RegisteredClaims
	UserID string `json:"user_id"`
	Nickname string `json:"nickname"`
	Email string `json:"email"`
}

func NewPayload(userID, nickname string, email string, ) *Payload {
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