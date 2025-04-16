package repositories

import (
	"context"

	"github.com/merdernoty/anime-service/internal/domain/models"
	"github.com/merdernoty/anime-service/internal/domain/repositories"
	"gorm.io/gorm"
)

func NewUserRepository(db *gorm.DB) repositories.UserRepository {
	return &UserRepositoryImpl{
		db: db,
	}
}

type UserRepositoryImpl struct {
	db *gorm.DB
}

func (r *UserRepositoryImpl) Create(ctx context.Context, user models.User) (models.User, error) {
	result := r.db.WithContext(ctx).Create(&user)
	if result.Error != nil {
		return models.User{}, result.Error
	}
	return user, nil
}

func (r *UserRepositoryImpl) GetByEmail(ctx context.Context, email string) (models.User, error) {
	var user models.User
	result := r.db.WithContext(ctx).Where("email = ?", email).First(&user)
	if result.Error != nil {
		return models.User{}, result.Error
	}
	return user, nil
}

func (r *UserRepositoryImpl) GetByID(ctx context.Context, id uint) (models.User, error) {
	var user models.User
	result := r.db.WithContext(ctx).First(&user, id)
	if result.Error != nil {
		return models.User{}, result.Error
	}
	return user, nil
}

func (r *UserRepositoryImpl) Update(ctx context.Context, user models.User) (models.User, error) {
	result := r.db.WithContext(ctx).Save(&user)
	if result.Error != nil {
		return models.User{}, result.Error
	}
	return user, nil
}

func (r *UserRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(context.Background()).Delete(&models.User{}, id).Error
}

func (r *UserRepositoryImpl) GetAll(ctx context.Context) ([]models.User, error) {
	var users []models.User
	result := r.db.WithContext(ctx).Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

func (r *UserRepositoryImpl) GetByNickName(ctx context.Context, nickName string) (models.User, error) {
	var user models.User
	result := r.db.WithContext(ctx).Where("nickname = ?", nickName).First(&user)
	if result.Error != nil {
		return models.User{}, result.Error
	}
	return user, nil
}