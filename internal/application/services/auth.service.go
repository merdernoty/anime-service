package services

import (
	"context"
	"emperror.dev/errors"
	"github.com/merdernoty/anime-service/internal/domain/dtos"
	"github.com/merdernoty/anime-service/internal/domain/models"
	"github.com/merdernoty/anime-service/internal/domain/repositories"
	"github.com/merdernoty/anime-service/pkg/auth"
	"logur.dev/logur"
)

var (
	ErrUserAlreadyExists  = errors.New("user with this email already exists")
	ErrUserNotFound       = errors.New("user not found")
	ErrUserNicknameExists = errors.New("user with this nickname already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type AuthServiceImpl struct {
	repo       repositories.UserRepository
	logger     logur.LoggerFacade
	tokenMaker auth.TokenMaker
}

func NewAuthService(repo repositories.UserRepository, logger logur.LoggerFacade, tokenMaker auth.TokenMaker) *AuthServiceImpl {
	return &AuthServiceImpl{
		repo:       repo,
		logger:     logger,
		tokenMaker: tokenMaker,
	}
}

func (s *AuthServiceImpl) Register(ctx context.Context, dto dtos.CreateUserDTO) (models.User, error) {
	_, err := s.repo.GetByEmail(ctx, dto.Email)
	if err != nil {
		return models.User{}, ErrUserAlreadyExists
	}

	_, err = s.repo.GetByNickName(ctx, dto.NickName)
	if err == nil {
		return models.User{}, ErrUserNicknameExists
	}

	user := models.User{
		Nickname:  dto.NickName,
		Email:     dto.Email,
		Firstname: dto.FirstName,
		Lastname:  dto.LastName,
		Password:  dto.Password,
	}

	if err := user.HashPassword(); err != nil {
		return models.User{}, errors.Wrap(err, "failed to hash password")
	}
	createdUser, err := s.repo.Create(ctx, user)
	if err != nil {
		return models.User{}, errors.Wrap(err, "failed to create user")
	}
	s.logger.Info("user registered successfully", map[string]interface{}{
		"NickName": user.Nickname,
		"Email":    user.Email,
	})
	createdUser.Password = ""

	return createdUser, nil
}

func (s *AuthServiceImpl) Login(ctx context.Context, dto dtos.LoginDTO) (dtos.TokenResponseDTO, error) {
	user, err := s.repo.GetByEmail(ctx, dto.Email)
	
	if err != nil {
		return dtos.TokenResponseDTO{}, ErrUserNotFound
	}

	if !user.CheckPassword(dto.Password) {
		return dtos.TokenResponseDTO{}, errors.New("invalid password")
	}

	token, err := s.tokenMaker.CreateToken(user.ID, user.Nickname, user.Email)

	if err != nil {
		return dtos.TokenResponseDTO{}, errors.Wrap(err, "failed to create token")
	}

	s.logger.Info("user logged in successfully", map[string]interface{}{
		"NickName": user.Nickname,
		"Email":    user.Email,
	})
	user.Password = ""
	return dtos.TokenResponseDTO{
		AccessToken: token,
	}, nil

}