package repositories

import (
	"context"

	"github.com/merdernoty/anime-service/internal/domain/models"
)

type UserRepository interface {
	Create(ctx context.Context, user models.User) (models.User, error)
	GetByEmail(ctx context.Context, email string) (models.User, error)
	GetByID(ctx context.Context, id uint) (models.User, error)
	Update(ctx context.Context, user models.User) (models.User, error)
	Delete(ctx context.Context, id uint) error
	GetAll(ctx context.Context) ([]models.User, error)
	GetByNickName(ctx context.Context, nickName string) (models.User, error)
	// GetUserFriends(ctx context.Context, userID uint) ([]models.User, error)
}