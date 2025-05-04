package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/merdernoty/anime-service/internal/interfaces/http/controllers"
    "github.com/merdernoty/anime-service/internal/interfaces/http/middleware"
)

type Service struct {
    AuthController *controllers.AuthController
    AnimeController *controllers.AnimeController
    UserController *controllers.UserController
}

func SetupRoutes(
    router *gin.Engine,
    service *Service,
    authMiddleware *middleware.AuthMiddleware,
) {
    api := router.Group("/api")
    RegisterUserRoutes(api, service.UserController, authMiddleware)
    RegisterAnimeRoutes(api, service.AnimeController, authMiddleware)
    RegisterAuthRoutes(api,service.AuthController)
}

func NewService(
    authController *controllers.AuthController,
    animeController *controllers.AnimeController,
    userController *controllers.UserController,
) *Service {
    return &Service{
        AuthController: authController,
        AnimeController: animeController,
        UserController: userController,
    }
}