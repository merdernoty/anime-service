package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/merdernoty/anime-service/internal/interfaces/http/controllers"
	"github.com/merdernoty/anime-service/internal/application/services"

)

func SetupAuthRoutes(router *gin.RouterGroup, authService *services.AuthServiceImpl) {
	authController := controllers.NewAuthController(authService)

	authRoutes := router.Group("/auth")
	{
		authRoutes.POST("/register", authController.Register)
		authRoutes.POST("/login", authController.Login)
	}
}	