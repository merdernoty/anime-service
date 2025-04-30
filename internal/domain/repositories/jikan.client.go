package repositories

import (
	"context"

	"github.com/merdernoty/anime-service/internal/domain/models"
)

type JikanClient interface {
	GetAnimeByID(ctx context.Context, malID int64) (*models.Anime, error)
	SearchAnime(ctx context.Context, query string, page, limit int) ([]*models.Anime, int, error)
	GetTopAnime(ctx context.Context, page, limit int) ([]*models.Anime, int, error)
	GetSeasonalAnime(ctx context.Context, year, season string, page, limit int) ([]*models.Anime, int, error)
	GetAnimeRecommendations(ctx context.Context, malID int64, page, limit int) ([]*models.Anime, error)
}