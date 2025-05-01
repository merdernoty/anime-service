package services

import (
	"context"
	"strconv"

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
	ErrTokenCreationFailed = errors.New("failed to create token")
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

func (s *AuthServiceImpl) Register(ctx context.Context, dto dtos.CreateUserDTO) (dtos.TokenResponseDTO, error) {
	_, err := s.repo.GetByEmail(ctx, dto.Email)
	if err == nil {
		return dtos.TokenResponseDTO{}, ErrUserAlreadyExists
	}

	_, err = s.repo.GetByNickName(ctx, dto.NickName)
	if err == nil {
		return dtos.TokenResponseDTO{}, ErrUserNicknameExists
	}

	user := models.User{
		Nickname:  dto.NickName,
		Email:     dto.Email,
		Firstname: dto.FirstName,
		Lastname:  dto.LastName,
		Password:  dto.Password,
	}

	if err := user.HashPassword(); err != nil {
		return dtos.TokenResponseDTO{}, errors.Wrap(err, "failed to hash password")
	}

	createdUser, err := s.repo.Create(ctx, user)
	if err != nil {
		return dtos.TokenResponseDTO{}, errors.Wrap(err, "failed to create user")
	}

	s.logger.Info("user registered successfully", map[string]interface{}{
		"NickName": user.Nickname,
		"Email":    user.Email,
	})

	createdUser.Password = ""

	accessToken, err := s.tokenMaker.CreateToken(
        user.ID,
        user.Nickname,
        user.Email,
        "access",
        3600, 
    )
	if err != nil {
		return dtos.TokenResponseDTO{}, errors.Wrap(err, "failed to create token")
	}

	refreshToken, err := s.tokenMaker.CreateToken(
        user.ID,
        user.Nickname,
        user.Email,
        "refresh",
        30*24*3600,
    )
    if err != nil {
        return dtos.TokenResponseDTO{}, errors.Wrap(err, "failed to create refresh token")
    }

	return dtos.TokenResponseDTO{
		AccessToken: accessToken,
		RefreshToken: refreshToken,
		TokenType:   "Bearer",
		ExpiresIn:    3600,
	}, nil
}
func (s *AuthServiceImpl) Login(ctx context.Context, dto dtos.LoginDTO) (dtos.TokenResponseDTO, error) {
    user, err := s.repo.GetByEmail(ctx, dto.Email)
    if err != nil {
        return dtos.TokenResponseDTO{}, ErrUserNotFound
    }

    if !user.CheckPassword(dto.Password) {
        return dtos.TokenResponseDTO{}, ErrInvalidCredentials
    }

    accessToken, err := s.tokenMaker.CreateToken(
        user.ID,
        user.Nickname,
        user.Email,
        "access",
        3600,
    )
    if err != nil {
        return dtos.TokenResponseDTO{}, ErrTokenCreationFailed
    }

    refreshToken, err := s.tokenMaker.CreateToken(
        user.ID,
        user.Nickname,
        user.Email,
        "refresh",
        30*24*3600,
    )
    if err != nil {
        return dtos.TokenResponseDTO{}, ErrTokenCreationFailed
    }

    s.logger.Info("user logged in successfully", map[string]interface{}{
        "NickName": user.Nickname,
        "Email":    user.Email,
    })

    return dtos.TokenResponseDTO{
        AccessToken:  accessToken,
        RefreshToken: refreshToken,
        TokenType:    "Bearer",
        ExpiresIn:    3600,
    }, nil
}

func (s *AuthServiceImpl) RefreshToken(ctx context.Context, dto dtos.RefreshTokenDTO) (string, error) {
	claims, err := s.tokenMaker.VerifyToken(dto.RefreshToken)
	if err != nil {
		return "", ErrInvalidCredentials
	}
	userID, err := strconv.ParseUint(claims.UserID, 10, 32)
	if err != nil {
		return "", errors.Wrap(err, "failed to parse user ID")
	}

	user, err := s.repo.GetByID(ctx, uint(userID))
	if err != nil {
		return "", ErrUserNotFound
	}

	token, err := s.tokenMaker.CreateToken(user.ID, user.Nickname, user.Email, "access", 3600)
	if err != nil {
		return "", ErrTokenCreationFailed
	}

	s.logger.Info("user refreshed token successfully", map[string]interface{}{
		"NickName": user.Nickname,
		"Email":    user.Email,
	})

	return token, nil
}
 