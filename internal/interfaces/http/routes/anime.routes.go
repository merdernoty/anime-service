package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/merdernoty/anime-service/internal/interfaces/http/controllers"
	"github.com/merdernoty/anime-service/internal/interfaces/http/middleware"
)

func RegisterAnimeRoutes(router *gin.RouterGroup, animeController *controllers.AnimeController, authMiddleware *middleware.AuthMiddleware) {
	publicAnime := router.Group("/anime")
	{
		publicAnime.GET("/:id", animeController.GetAnimeByID)
		publicAnime.GET("/search", animeController.SearchAnime)
		publicAnime.GET("/top", animeController.GetTopAnime)
		publicAnime.GET("/seasonal/:year/:season", animeController.GetSeasonalAnime)
		publicAnime.GET("/:id/recommendations", animeController.GetAnimeRecommendations)
	}

	authorizedUser := router.Group("/me")
	authorizedUser.Use(authMiddleware.Auth())
	{
		myAnime := authorizedUser.Group("/anime")
		{
			myAnime.GET("", animeController.GetUserAnimeList)
			myAnime.POST("", animeController.AddAnimeToUserList)
			myAnime.DELETE("/:anime_id", animeController.RemoveAnimeFromUserList)
			myAnime.PUT("/:anime_id/status", animeController.UpdateUserAnimeStatus)
			myAnime.PUT("/:anime_id/episodes", animeController.UpdateUserAnimeEpisodes)
			myAnime.PUT("/:anime_id/rating", animeController.UpdateUserAnimeRating)
			myAnime.GET("/stats", animeController.GetUserAnimeStats)
		}
	}
	
	adminRoutes := router.Group("/users")
	adminRoutes.Use(authMiddleware.Auth())
	{
		adminRoutes.GET("/:user_id/anime", animeController.GetUserAnimeList)
		adminRoutes.POST("/:user_id/anime", animeController.AddAnimeToUserList)
		adminRoutes.DELETE("/:user_id/anime/:anime_id", animeController.RemoveAnimeFromUserList)
		adminRoutes.PUT("/:user_id/anime/:anime_id/status", animeController.UpdateUserAnimeStatus)
		adminRoutes.PUT("/:user_id/anime/:anime_id/episodes", animeController.UpdateUserAnimeEpisodes)
		adminRoutes.PUT("/:user_id/anime/:anime_id/rating", animeController.UpdateUserAnimeRating)
		adminRoutes.GET("/:user_id/anime/stats", animeController.GetUserAnimeStats)
	}
}