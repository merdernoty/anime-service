package services

import (
	"emperror.dev/errors"
	"github.com/gin-gonic/gin"
	"github.com/merdernoty/anime-service/internal/domain/dtos"
	"github.com/merdernoty/anime-service/internal/domain/repositories"
	"logur.dev/logur"
)

type UserServiceImlp struct {
	repo repositories.UserRepository
	logger         logur.LoggerFacade
}


func NewUserService(repo repositories.UserRepository, logger logur.LoggerFacade) *UserServiceImlp {
	return &UserServiceImlp{
		repo:   repo,
		logger: logger,
	}
}

func (s *UserServiceImlp) GetUserProfile(ctx *gin.Context, id uint) (dtos.UserResponseDTO, error) {
	user, err:= s.repo.GetByID(ctx, id)
	if err != nil {
		return dtos.UserResponseDTO{}, ErrUserNotFound
	}
	user.Password = ""
	return dtos.UserResponseDTO{
		ID:        user.ID,
		Email:     user.Email,
		NickName:  user.Nickname,
		FirstName: user.Firstname,
		LastName:  user.Lastname,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (s *UserServiceImlp) UpdateUserProfile(ctx *gin.Context, id uint, dto dtos.UpdateUserDTO) (dtos.UserResponseDTO, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return dtos.UserResponseDTO{}, ErrUserNotFound
	}

	if dto.NickName != "" {
		if *&dto.NickName != user.Nickname {
			existingUser, err := s.repo.GetByNickName(ctx, *&dto.NickName)
			if err == nil && existingUser.ID != user.ID {
				return dtos.UserResponseDTO{}, errors.New("nickname already exists")
			}
			user.Nickname = *&dto.NickName
		}
	}

	if dto.Email != "" {
		if *&dto.Email != user.Email {
			existingUser, err := s.repo.GetByEmail(ctx, *&dto.Email)
			if err == nil && existingUser.ID != user.ID {
				return dtos.UserResponseDTO{}, errors.New("пользователь с таким email уже существует")
			}
			user.Email = *&dto.Email
		}
	}

	if dto.FirstName != "" {
		user.Firstname = *&dto.FirstName
	}

	if dto.LastName != "" {
		user.Lastname = *&dto.LastName
	}


    updatedUser, err := s.repo.Update(ctx, user)
	if err != nil {
		return dtos.UserResponseDTO{}, err
	} 

	return dtos.ToUserResponse(updatedUser), nil
}
