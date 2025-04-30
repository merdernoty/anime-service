package services

import (
	"context"

	"github.com/merdernoty/anime-service/internal/domain/models"
)

type AnimeService interface {
	GetAnimeByID(ctx context.Context, malID int64) (*models.Anime, error)
	SearchAnime(ctx context.Context, query string, page, limit int) ([]*models.Anime, int, error)
	GetTopAnime(ctx context.Context, page, limit int) ([]*models.Anime, int, error)
	GetSeasonalAnime(ctx context.Context, year, season string, page, limit int) ([]*models.Anime, int, error)
	GetAnimeRecommendations(ctx context.Context, malID int64, page, limit int) ([]*models.Anime, error)

	GetUserAnimeList(ctx context.Context, filter models.UserAnimeFilter) (*models.UserAnimeList, error)
	AddAnimeToUserList(ctx context.Context, userID uint, animeMALID int64, status models.WatchStatus) error
	RemoveAnimeFromUserList(ctx context.Context, userID uint, animeMALID int64) error
	UpdateAnimeStatus(ctx context.Context, userID uint, animeMALID int64, status models.WatchStatus) error
	UpdateUserAnimeEpisodes(ctx context.Context, userID uint, animeMALID int64, episodesWatched int) error
	UpdateUserAnimeRating(ctx context.Context, userID uint, animeMALID int64, rating float32) error
	GetUserAnimeStats(ctx context.Context, userID uint) (*models.AnimeStats, error)
}