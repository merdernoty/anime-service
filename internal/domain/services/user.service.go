package services

import (
	"context"

	"github.com/merdernoty/anime-service/internal/domain/models"
)

type UserService interface {
	GetUserProfile(ctx context.Context, userID uint) (*models.User, error)
	UpdateUserProfile(ctx context.Context, userID uint, profile *models.User) error
	GetUserFriends(ctx context.Context, userID uint) ([]*models.User, error)
}