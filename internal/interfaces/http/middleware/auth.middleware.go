package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/merdernoty/anime-service/internal/domain/dtos"
	"github.com/merdernoty/anime-service/internal/domain/repositories"
	"github.com/merdernoty/anime-service/pkg/auth"
	"strconv"
)

type AuthMiddleware struct {
	tokenMaker auth.TokenMaker
	userRepository repositories.UserRepository
}


func NewAuthMiddleware(tokenMaker auth.TokenMaker, userRepository repositories.UserRepository) *AuthMiddleware {
	return &AuthMiddleware{
		tokenMaker: tokenMaker,
		userRepository: userRepository,
	}
}

func (m *AuthMiddleware) Auth() gin.HandlerFunc {
	return func (ctx *gin.Context) {
		token := ctx.GetHeader("Authorization")
		if token == "" {
			ctx.JSON(http.StatusUnauthorized, dtos.ErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: "Authorization header is missing",
			})
			ctx.Abort()
			return
		}

		parts := strings.Split(token, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			ctx.JSON(http.StatusUnauthorized, dtos.ErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: "Invalid authorization header format",
			})
			ctx.Abort()
			return
		}

		accessToken := parts[1]

		payload, err := m.tokenMaker.VerifyToken(accessToken)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, dtos.ErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: "Invalid or expired token",
			})
			ctx.Abort()
			return
		}

		userID, err := strconv.ParseUint(payload.UserID, 10, 32)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, dtos.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "Invalid user ID format",
			})
			ctx.Abort()
			return
		}

		user, err := m.userRepository.GetByID(ctx, uint(userID))
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, dtos.ErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: "User not found",
			})
			ctx.Abort()
			return
		}


		ctx.Set("user", user)
		ctx.Set("payload", payload)
		ctx.Next()
	}
}