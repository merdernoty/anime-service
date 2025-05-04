package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/merdernoty/anime-service/internal/interfaces/http/controllers"
	"github.com/merdernoty/anime-service/internal/interfaces/http/middleware"
)

func RegisterUserRoutes(router *gin.RouterGroup, userController *controllers.UserController, authMiddleware *middleware.AuthMiddleware) {
		userRoutes := router.Group("/user")
		userRoutes.Use(authMiddleware.Auth())
		{	
			userRoutes.GET("/profile", func(ctx *gin.Context) {
				response, err := userController.GetUserProfile(ctx)
				if err != nil {
					ctx.JSON(500, gin.H{"error": err.Error()})
					return
				}
				ctx.JSON(200, response)
			})
			userRoutes.PUT("/profile", func(ctx *gin.Context) {
				response, err := userController.UpdateUserProfile(ctx)
				if err != nil {
					ctx.JSON(500, gin.H{"error": err.Error()})
					return
				}
				ctx.JSON(200, response)
			})
		}
}