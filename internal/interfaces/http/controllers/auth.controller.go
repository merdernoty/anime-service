package controllers

import (
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

func handleError(ctx *gin.Context, err error) {
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

func (c *AuthController) Register(ctx *gin.Context) {
	var dto dtos.CreateUserDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "validation error", "details": err.Error()})
		return
	}

	user, err := c.authService.Register(ctx, dto)
	if err != nil {
		handleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, dtos.ToUserResponse(user))
}

func (c *AuthController) Login(ctx *gin.Context) {
	var dto dtos.LoginDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "validation error", "details": err.Error()})
		return
	}

	tokenResponse, err := c.authService.Login(ctx, dto)
	if err != nil {
		handleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, tokenResponse)
}
