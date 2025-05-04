package controllers

import (
	"net/http"

	"emperror.dev/errors"
	"github.com/gin-gonic/gin"
	"github.com/merdernoty/anime-service/internal/application/services"
	"github.com/merdernoty/anime-service/internal/domain/dtos"
	"logur.dev/logur"
)

type UserController struct {
	UserService *services.UserServiceImlp
	logger logur.LoggerFacade
}

func NewUserController(userService *services.UserServiceImlp, logger logur.LoggerFacade) *UserController {
	return &UserController{
		UserService: userService,
		logger: logger,
	}
}

func handleUserError(ctx *gin.Context, err error){
	switch {
		case err == services.ErrUserNotFound:
			ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		case err == services.ErrUserAlreadyExists:
			ctx.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		case err == services.ErrInvalidCredentials:
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		case err == services.ErrUnauthorized:
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
	}
}

// GetUserProfile godoc
// @Summary      Получить профиль пользователя
// @Description  Получить профиль пользователя по ID
// @Tags         Profile
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  dtos.UserResponseDTO  "Успешное получение профиля"
// @Failure      400  {object}  dtos.ErrorResponse    "Ошибка валидации"
// @Failure      401  {object}  dtos.ErrorResponse    "Пользователь не авторизован"
// @Failure      404  {object}  dtos.ErrorResponse    "Пользователь не найден"
// @Failure      500  {object}  dtos.ErrorResponse    "Внутренняя ошибка сервера"
// @Router       /user/profile [get]
func (c *UserController) GetUserProfile(ctx *gin.Context) (dtos.UserResponseDTO, error) {
    userIDRaw, exists := ctx.Get("userID")
    if !exists {
        return dtos.UserResponseDTO{}, services.ErrUnauthorized
    }
	userID, ok := userIDRaw.(uint)
	if !ok {
		c.logger.Error("Failed to cast userID to uint")
		return dtos.UserResponseDTO{}, errors.New("failed to cast userID to uint")
	}
	profile, err := c.UserService.GetUserProfile(ctx, userID)
    if err != nil {
        return dtos.UserResponseDTO{}, err
    }

    userDTO := dtos.UserResponseDTO{
        ID:        profile.ID,
        Email:     profile.Email,
        NickName:  profile.NickName,
        FirstName: profile.FirstName,
        LastName:  profile.LastName,
        CreatedAt: profile.CreatedAt,
        UpdatedAt: profile.UpdatedAt,
    }
    
    return userDTO, nil
}
// UpdateUserProfile godoc
// @Summary		Обновить профиль пользователя
// @Description	Обновить профиль пользователя по ID
// @Tags			Profile
// @Accept			json
// @Produce		json
// @Security     BearerAuth
// @Param request body dtos.UpdateUserDTO false "Данные для обновления профиля"
// @Success		200		{object}	dtos.UserResponseDTO	"Успешное обновление профиля"
// @Failure		400		{object}	dtos.ErrorResponse		"Ошибка валидации"
// @Failure		404		{object}	dtos.ErrorResponse		"Пользователь не найден"
// @Failure		500		{object}	dtos.ErrorResponse		"Внутренняя ошибка сервера"
// @Router			/user/profile [put]
func (c *UserController) UpdateUserProfile(ctx *gin.Context) (dtos.UserResponseDTO, error) {
	userIDRaw, exists := ctx.Get("userID")
    if !exists {
        return dtos.UserResponseDTO{}, services.ErrUnauthorized
    }
	userID, ok := userIDRaw.(uint)
	if !ok {
		c.logger.Error("Failed to cast userID to uint")
		return dtos.UserResponseDTO{}, errors.New("failed to cast userID to uint")
	}
	var dto dtos.UpdateUserDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "validation error", "details": err.Error()})
		return dtos.UserResponseDTO{}, err
	}
	profile, err := c.UserService.UpdateUserProfile(ctx, userID, dto)
	if err != nil {	
		handleUserError(ctx, err)
		return dtos.UserResponseDTO{}, services.ErrUserNotFound
	}
	return dtos.UserResponseDTO{
		ID:        profile.ID,
		Email:     profile.Email,
		NickName:  profile.NickName,
		FirstName: profile.FirstName,
		LastName:  profile.LastName,
		CreatedAt: profile.CreatedAt,
		UpdatedAt: profile.UpdatedAt,
	}, nil
}

