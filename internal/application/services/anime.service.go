package services

import (
	"context"
	"time"

	"emperror.dev/errors"
	"github.com/merdernoty/anime-service/internal/domain/models"
	"github.com/merdernoty/anime-service/internal/infrastructure/api"
	"github.com/merdernoty/anime-service/internal/infrastructure/repositories"
	"logur.dev/logur"
)

type AnimeServiceImpl struct {
	jikanClient   *api.JikanClient
	userAnimeRepo *repositories.UserAnimeRepository
	logger        logur.LoggerFacade
}

var (
	ErrAnimeNotFound = errors.New("anime not found")
	ErrAnimeAlreadyExists = errors.New("anime already exists in user list")
	ErrAnimeNotInUserList = errors.New("anime not found in user list")
	ErrAnimeUpdateFailed = errors.New("failed to update anime in user list")
	ErrAnimeDeleteFailed = errors.New("failed to delete anime from user list")
	ErrAnimeStatsFailed = errors.New("failed to get user anime stats")
	ErrFetchAnimeFailed = errors.New("failed to fetch anime details")
)

func NewAnimeService( jikanClient *api.JikanClient, userAnimeRepo *repositories.UserAnimeRepository, logger logur.LoggerFacade) *AnimeServiceImpl {
	return &AnimeServiceImpl{
		jikanClient:   jikanClient,
		userAnimeRepo: userAnimeRepo,
		logger:        logger,
	}
}
func (s *AnimeServiceImpl) GetAnimeByID(ctx context.Context, malID int64) (*models.Anime, error) {
	s.logger.Info("Getting anime by ID", map[string]interface{}{
		"mal_id": malID,
	})

	anime, err := s.jikanClient.GetAnimeByID(ctx, malID)
	if err != nil {
		s.logger.Error("Error getting anime by ID", map[string]interface{}{
			"mal_id": malID,
			"error":  err.Error(),
		})
		return nil, ErrAnimeNotFound
	}

	return anime, nil
}

func (s *AnimeServiceImpl) SearchAnime(ctx context.Context, query string, page, limit int) ([]*models.Anime, int, error) {
	s.logger.Info("Searching anime", map[string]interface{}{
		"query": query,
		"page":  page,
		"limit": limit,
	})

	animes, totalPages, err := s.jikanClient.SearchAnime(ctx, query, page, limit)
	if err != nil {
		s.logger.Error("Error searching anime", map[string]interface{}{
			"query": query,
			"error": err.Error(),
		})
		return nil, 0, ErrFetchAnimeFailed
	}

	return animes, totalPages, nil
}

func (s *AnimeServiceImpl) GetTopAnime(ctx context.Context, page, limit int) ([]*models.Anime, int, error) {
	s.logger.Info("Getting top anime", map[string]interface{}{
		"page":  page,
		"limit": limit,
	})

	animes, totalPages, err := s.jikanClient.GetTopAnime(ctx, page, limit)
	if err != nil {
		s.logger.Error("Error getting top anime", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, 0, ErrFetchAnimeFailed
	}

	return animes, totalPages, nil
}

func (s *AnimeServiceImpl) GetSeasonalAnime(ctx context.Context, year, season string, page, limit int) ([]*models.Anime, int, error) {
	s.logger.Info("Getting seasonal anime", map[string]interface{}{
		"year":   year,
		"season": season,
		"page":   page,
		"limit":  limit,
	})

	animes, totalPages, err := s.jikanClient.GetSeasonalAnime(ctx, year, season, page, limit)
	if err != nil {
		s.logger.Error("Error getting seasonal anime", map[string]interface{}{
			"year":   year,
			"season": season,
			"error":  err.Error(),
		})
		return nil, 0, ErrFetchAnimeFailed
	}

	return animes, totalPages, nil
}

func (s *AnimeServiceImpl) GetAnimeRecommendations(ctx context.Context, malID int64, page, limit int) ([]*models.Anime, error) {
	s.logger.Info("Getting anime recommendations", map[string]interface{}{
		"mal_id": malID,
		"page":   page,
		"limit":  limit,
	})

	animes, err := s.jikanClient.GetAnimeRecommendations(ctx, malID, page, limit)
	if err != nil {
		s.logger.Error("Error getting anime recommendations", map[string]interface{}{
			"mal_id": malID,
			"error":  err.Error(),
		})
		return nil, ErrFetchAnimeFailed
	}

	return animes, nil
}

func (s *AnimeServiceImpl) GetUserAnimeList(ctx context.Context, filter models.UserAnimeFilter) (*models.UserAnimeList, error) {
	s.logger.Info("Getting user anime list", map[string]interface{}{
		"user_id": filter.UserID,
		"status":  filter.Status,
		"page":    filter.Page,
		"limit":   filter.Limit,
	})

	userAnimeList, err := s.userAnimeRepo.GetUserAnimeWithDetails(ctx, filter, s.jikanClient)
	if err != nil {
		s.logger.Error("Error getting user anime list", map[string]interface{}{
			"user_id": filter.UserID,
			"error":   err.Error(),
		})
		return nil, ErrFetchAnimeFailed
	}

	return userAnimeList, nil
}

func (s *AnimeServiceImpl) AddAnimeToUserList(ctx context.Context, userID uint, animeMALID int64, status models.WatchStatus) error {
	s.logger.Info("Adding anime to user list", map[string]interface{}{
		"user_id":      userID,
		"anime_mal_id": animeMALID,
		"status":       status,
	})

	_, err := s.jikanClient.GetAnimeByID(ctx, animeMALID)
	if err != nil {
		s.logger.Error("Error getting anime by ID for adding to user list", map[string]interface{}{
			"anime_mal_id": animeMALID,
			"error":        err.Error(),
		})
		return ErrAnimeNotFound
	}

	userAnime := &models.UserAnime{
		UserID:     userID,
		AnimeMALID: animeMALID,
		Status:     status,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	err = s.userAnimeRepo.CreateOrUpdateUserAnime(ctx, userAnime)
	if err != nil {
		s.logger.Error("Error adding anime to user list", map[string]interface{}{
			"user_id":      userID,
			"anime_mal_id": animeMALID,
			"error":        err.Error(),
		})
		return ErrAnimeUpdateFailed
	}

	return nil
}

func (s *AnimeServiceImpl) RemoveAnimeFromUserList(ctx context.Context, userID uint, animeMALID int64) error {
	s.logger.Info("Removing anime from user list", map[string]interface{}{
		"user_id":      userID,
		"anime_mal_id": animeMALID,
	})

	err := s.userAnimeRepo.DeleteByUserAndAnimeMALID(ctx, userID, animeMALID)
	if err != nil {
		s.logger.Error("Error removing anime from user list", map[string]interface{}{
			"user_id":      userID,
			"anime_mal_id": animeMALID,
			"error":        err.Error(),
		})
		return ErrAnimeDeleteFailed
	}

	return nil
}

func (s *AnimeServiceImpl) UpdateUserAnimeStatus(ctx context.Context, userID uint, animeMALID int64, status models.WatchStatus) error {
	s.logger.Info("Updating user anime status", map[string]interface{}{
		"user_id":      userID,
		"anime_mal_id": animeMALID,
		"status":       status,
	})

	err := s.userAnimeRepo.ChangeStatus(ctx, userID, animeMALID, status)
	if err != nil {
		s.logger.Error("Error updating user anime status", map[string]interface{}{
			"user_id":      userID,
			"anime_mal_id": animeMALID,
			"error":        err.Error(),
		})
		return ErrAnimeUpdateFailed
	}

	return nil
}

func (s *AnimeServiceImpl) UpdateUserAnimeEpisodes(ctx context.Context, userID uint, animeMALID int64, episodesWatched int) error {
	s.logger.Info("Updating user anime episodes watched", map[string]interface{}{
		"user_id":          userID,
		"anime_mal_id":     animeMALID,
		"episodes_watched": episodesWatched,
	})

	err := s.userAnimeRepo.UpdateEpisodesWatched(ctx, userID, animeMALID, episodesWatched)
	if err != nil {
		s.logger.Error("Error updating user anime episodes watched", map[string]interface{}{
			"user_id":      userID,
			"anime_mal_id": animeMALID,
			"error":        err.Error(),
		})
		return ErrAnimeUpdateFailed
	}

	return nil
}

func (s *AnimeServiceImpl) UpdateUserAnimeRating(ctx context.Context, userID uint, animeMALID int64, rating float32) error {
	s.logger.Info("Updating user anime rating", map[string]interface{}{
		"user_id":      userID,
		"anime_mal_id": animeMALID,
		"rating":       rating,
	})

	err := s.userAnimeRepo.UpdateRating(ctx, userID, animeMALID, rating)
	if err != nil {
		s.logger.Error("Error updating user anime rating", map[string]interface{}{
			"user_id":      userID,
			"anime_mal_id": animeMALID,
			"error":        err.Error(),
		})
		return ErrAnimeUpdateFailed
	}

	return nil
}

func (s *AnimeServiceImpl) GetUserAnimeStats(ctx context.Context, userID uint) (*models.AnimeStats, error) {
	s.logger.Info("Getting user anime stats", map[string]interface{}{
		"user_id": userID,
	})

	stats, err := s.userAnimeRepo.GetUserStats(ctx, userID)
	if err != nil {
		s.logger.Error("Error getting user anime stats", map[string]interface{}{
			"user_id": userID,
			"error":   err.Error(),
		})
		return nil, ErrAnimeStatsFailed
	}

	return stats, nil
}