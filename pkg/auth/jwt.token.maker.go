package auth

import (
	"strconv"
	"time"

	"emperror.dev/errors"
	"github.com/golang-jwt/jwt/v5"
)

type JWTTokenMaker struct {
	secretKey string
	duration time.Duration
}

func NewJWTTokenMaker(secretKey string, duration time.Duration) *JWTTokenMaker {
	return &JWTTokenMaker{
		secretKey: secretKey,
		duration:  duration,
	}
}

func (maker *JWTTokenMaker) CreateToken(
    userID uint,
    nickname string, 
	email string,
    token_type string,
    expiresIn time.Duration,
) (string, error) {
    payload := &Payload{
        UserID:    strconv.FormatUint(uint64(userID), 10),
        Nickname:  nickname,
        Email:     email,
        TokenType: token_type,
        IssuedAt:  time.Now().Unix(),
        ExpiresAt: time.Now().Add(expiresIn).Unix(), 
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
    return token.SignedString([]byte(maker.secretKey))
}
func (maker *JWTTokenMaker) VerifyToken(tokenString string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(maker.secretKey), nil
	}

	token, err := jwt.ParseWithClaims(
		tokenString, 
		&Payload{}, 
		keyFunc,
	)
	if err != nil {
		switch {
		case err.Error() == "token is expired":
			return nil, errors.New("token is expired")
		case err.Error() == "token is malformed":
			return nil, errors.New("token is malformed")
		case err.Error() == "token is unverifiable":
			return nil, errors.New("token is unverifiable")
		default:
			return nil, errors.Wrap(err, "failed to parse token")
		}
	}

	payload, ok := token.Claims.(*Payload)
	if !ok {
		return nil, errors.New("failed to parse token claims")
	}

	return payload, nil
}