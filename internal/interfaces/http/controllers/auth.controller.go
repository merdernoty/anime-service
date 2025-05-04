package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/merdernoty/anime-service/internal/application/services"
	"github.com/merdernoty/anime-service/internal/domain/dtos"
)

type AuthController struct {
	authService *services.AuthServiceImpl
}

func NewAuthController(authService *services.AuthServiceImpl) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

func handleAuthError(ctx *gin.Context, err error) {
	switch {
	case err == services.ErrUserAlreadyExists:
		ctx.JSON(http.StatusConflict, gin.H{
			"error":   "user already exist",
			"details": err.Error(),
		})
	case err == services.ErrInvalidCredentials:
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error":   "invalid credentials",
			"details": err.Error(),
		})
	default:
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "internal server error",
			"details": err.Error(),
		})
	}
}


// Register godoc
//	@Summary		Регистрация нового пользователя
//	@Description	Создает нового пользователя в системе
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dtos.CreateUserDTO		true	"Данные для регистрации"
//	@Success		201		{object}	dtos.UserResponseDTO	"Успешная регистрация"
//	@Failure		400		{object}	dtos.ErrorResponse		"Ошибка валидации"
//	@Failure		409		{object}	dtos.ErrorResponse		"Пользователь уже существует"
//	@Failure		500		{object}	dtos.ErrorResponse		"Внутренняя ошибка сервера"
//	@Router			/auth/register [post]
func (c *AuthController) Register(ctx *gin.Context) {
	var dto dtos.CreateUserDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "validation error", "details": err.Error()})
		return
	}

	tokenResponse, err := c.authService.Register(ctx, dto, ctx.Writer)
	if err != nil {
		fmt.Print(err)
		handleAuthError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, tokenResponse)
}


// Login godoc
//	@Summary		Аутентификация пользователя
//	@Description	Выполняет вход в систему и возвращает токены
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dtos.LoginDTO			true	"Учетные данные для входа"
//	@Success		200		{object}	dtos.TokenResponseDTO	"Успешная аутентификация"
//	@Failure		400		{object}	dtos.ErrorResponse		"Ошибка валидации"
//	@Failure		401		{object}	dtos.ErrorResponse		"Неверные учетные данные"
//	@Failure		500		{object}	dtos.ErrorResponse		"Внутренняя ошибка сервера"
//	@Router			/auth/login [post]
func (c *AuthController) Login(ctx *gin.Context) {
	var dto dtos.LoginDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "validation error", "details": err.Error()})
		return
	}

	tokenResponse, err := c.authService.Login(ctx, dto, ctx.Writer)
	if err != nil {
		handleAuthError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, tokenResponse)
}


// RefreshToken обновляет access-токен
//	@Summary		Обновление JWT-токена
//	@Description	Обновляет access-токен, используя refresh-токен из куки
//	@Tags			Auth
//	@Produce		json
//	@Success		200	{object}	dtos.TokenResponseDTO	"Возвращает новый access-токен"
//	@Failure		401	{object}	dtos.ErrorResponse		"Refresh-токен недействителен или отсутствует"
//	@Router			/auth/refresh [post]
func (c *AuthController) RefreshToken(ctx *gin.Context) {
	tokenResponse, err := c.authService.RefreshToken(ctx, ctx.Request, ctx.Writer)
	if err != nil {
		handleAuthError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, tokenResponse)
}

// Logout выход из системы
//	@Summary		Выход из системы
//	@Description	Удаляет refresh-токен и очищает куки
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	map[string]string	"Успешный выход"
//	@Router			/auth/logout [post]
func (c *AuthController) Logout(ctx *gin.Context) {
	c.authService.Logout(ctx)
	ctx.JSON(http.StatusOK, gin.H{
		"message": "successfully logged out",
	})
}