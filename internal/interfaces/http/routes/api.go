package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/merdernoty/anime-service/internal/interfaces/http/controllers"
)

type Service struct {
    AuthController *controllers.AuthController
}

func SetupRoutes(
    router *gin.Engine,
    service *Service,
) {
    api := router.Group("/api")

    public := api.Group("")
    {
        public.POST("/auth/register", service.AuthController.Register)
        public.POST("/auth/login", service.AuthController.Login)
    }

    // protected := api.Group("")
    // protected.Use(middleware.AuthMiddleware(service.AuthService))
    {
    // protected.GET("/profile", profileController.GetProfile)
    }
}

func NewService(
    authController *controllers.AuthController,
) *Service {
    return &Service{
        AuthController: authController,
    }
}