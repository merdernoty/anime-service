package services

import (
	"context"

	"github.com/merdernoty/anime-service/internal/domain/dtos"
	"github.com/merdernoty/anime-service/internal/domain/models"
)
type AuthService interface {
	Register(ctx context.Context, dto dtos.CreateUserDTO) (models.User, error)
	Login(ctx context.Context, dto dtos.LoginDTO) (models.User, error)
	RefreshToken(ctx context.Context, dto dtos.RefreshTokenDTO) (string, error)
}