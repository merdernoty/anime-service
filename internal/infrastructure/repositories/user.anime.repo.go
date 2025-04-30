package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"emperror.dev/errors"
	"github.com/merdernoty/anime-service/internal/domain/models"
	"github.com/merdernoty/anime-service/internal/infrastructure/api"
	"logur.dev/logur"
)

type UserAnimeRepository struct {
	db     *sql.DB
	logger logur.LoggerFacade
}

func NewUserAnimeRepository(db *sql.DB, logger logur.LoggerFacade) *UserAnimeRepository {
	return &UserAnimeRepository{
		db:     db,
		logger: logger,
	}
}

func (r *UserAnimeRepository) GetByID(ctx context.Context, id uint) (*models.UserAnime, error) {
	query := `
		SELECT id, user_id, anime_mal_id, status, rating, notes, episodes_watched, created_at, updated_at
		FROM user_animes
		WHERE id = $1
	`

	row := r.db.QueryRowContext(ctx, query, id)
	userAnime := &models.UserAnime{}

	err := row.Scan(
		&userAnime.ID,
		&userAnime.UserID,
		&userAnime.AnimeMALID,
		&userAnime.Status,
		&userAnime.Rating,
		&userAnime.Notes,
		&userAnime.EpisodesWatched,
		&userAnime.CreatedAt,
		&userAnime.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("user anime not found")
	}

	if err != nil {
		r.logger.Error("Error getting user anime by ID", map[string]interface{}{
			"id":    id,
			"error": err.Error(),
		})
		return nil, errors.Wrap(err, "error getting user anime by ID")
	}

	return userAnime, nil
}

func (r *UserAnimeRepository) GetByUserAndAnimeMALID(ctx context.Context, userID uint, animeMALID int64) (*models.UserAnime, error) {
	query := `
		SELECT id, user_id, anime_mal_id, status, rating, notes, episodes_watched, created_at, updated_at
		FROM user_animes
		WHERE user_id = $1 AND anime_mal_id = $2
	`

	row := r.db.QueryRowContext(ctx, query, userID, animeMALID)
	userAnime := &models.UserAnime{}

	err := row.Scan(
		&userAnime.ID,
		&userAnime.UserID,
		&userAnime.AnimeMALID,
		&userAnime.Status,
		&userAnime.Rating,
		&userAnime.Notes,
		&userAnime.EpisodesWatched,
		&userAnime.CreatedAt,
		&userAnime.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("user anime not found")
	}

	if err != nil {
		r.logger.Error("Error getting user anime by user ID and anime MAL ID", map[string]interface{}{
			"user_id":      userID,
			"anime_mal_id": animeMALID,
			"error":        err.Error(),
		})
		return nil, errors.Wrap(err, "error getting user anime by user ID and anime MAL ID")
	}

	return userAnime, nil
}

func (r *UserAnimeRepository) List(ctx context.Context, filter models.UserAnimeFilter) ([]*models.UserAnime, int, error) {
	var conditions []string
	var args []interface{}
	argCounter := 1

	conditions = append(conditions, fmt.Sprintf("user_id = $%d", argCounter))
	args = append(args, filter.UserID)
	argCounter++

	if filter.Status != "" {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argCounter))
		args = append(args, filter.Status)
		argCounter++
	}

	whereClause := strings.Join(conditions, " AND ")

	countQuery := fmt.Sprintf(`
		SELECT COUNT(*) FROM user_animes WHERE %s
	`, whereClause)

	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		r.logger.Error("Error counting user animes", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, 0, errors.Wrap(err, "error counting user animes")
	}

	limit := 10
	if filter.Limit > 0 {
		limit = filter.Limit
	}

	offset := 0
	if filter.Page > 1 {
		offset = (filter.Page - 1) * limit
	}

	query := fmt.Sprintf(`
		SELECT id, user_id, anime_mal_id, status, rating, notes, episodes_watched, created_at, updated_at
		FROM user_animes
		WHERE %s
		ORDER BY updated_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argCounter, argCounter+1)

	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		r.logger.Error("Error listing user animes", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, 0, errors.Wrap(err, "error listing user animes")
	}
	defer rows.Close()

	var userAnimes []*models.UserAnime
	for rows.Next() {
		userAnime := &models.UserAnime{}
		err := rows.Scan(
			&userAnime.ID,
			&userAnime.UserID,
			&userAnime.AnimeMALID,
			&userAnime.Status,
			&userAnime.Rating,
			&userAnime.Notes,
			&userAnime.EpisodesWatched,
			&userAnime.CreatedAt,
			&userAnime.UpdatedAt,
		)
		if err != nil {
			r.logger.Error("Error scanning user anime row", map[string]interface{}{
				"error": err.Error(),
			})
			return nil, 0, errors.Wrap(err, "error scanning user anime row")
		}
		userAnimes = append(userAnimes, userAnime)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, errors.Wrap(err, "error iterating user anime rows")
	}

	return userAnimes, total, nil
}

func (r *UserAnimeRepository) Create(ctx context.Context, userAnime *models.UserAnime) error {
	query := `
		INSERT INTO user_animes (
			user_id, anime_mal_id, status, rating, notes, episodes_watched, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8
		) RETURNING id
	`

	now := time.Now()
	userAnime.CreatedAt = now
	userAnime.UpdatedAt = now

	err := r.db.QueryRowContext(ctx, query,
		userAnime.UserID,
		userAnime.AnimeMALID,
		userAnime.Status,
		userAnime.Rating,
		userAnime.Notes,
		userAnime.EpisodesWatched,
		userAnime.CreatedAt,
		userAnime.UpdatedAt,
	).Scan(&userAnime.ID)

	if err != nil {
		r.logger.Error("Error creating user anime", map[string]interface{}{
			"user_id":      userAnime.UserID,
			"anime_mal_id": userAnime.AnimeMALID,
			"error":        err.Error(),
		})
		return errors.Wrap(err, "error creating user anime")
	}

	return nil
}

func (r *UserAnimeRepository) Update(ctx context.Context, userAnime *models.UserAnime) error {
	query := `
		UPDATE user_animes SET
			status = $1, 
			rating = $2, 
			notes = $3, 
			episodes_watched = $4,
			updated_at = $5
		WHERE id = $6
	`

	userAnime.UpdatedAt = time.Now()

	result, err := r.db.ExecContext(ctx, query,
		userAnime.Status,
		userAnime.Rating,
		userAnime.Notes,
		userAnime.EpisodesWatched,
		userAnime.UpdatedAt,
		userAnime.ID,
	)

	if err != nil {
		r.logger.Error("Error updating user anime", map[string]interface{}{
			"id":    userAnime.ID,
			"error": err.Error(),
		})
		return errors.Wrap(err, "error updating user anime")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "error getting rows affected")
	}

	if rowsAffected == 0 {
		return errors.New("user anime not found")
	}

	return nil
}

func (r *UserAnimeRepository) Delete(ctx context.Context, id uint) error {
	query := `DELETE FROM user_animes WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		r.logger.Error("Error deleting user anime", map[string]interface{}{
			"id":    id,
			"error": err.Error(),
		})
		return errors.Wrap(err, "error deleting user anime")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "error getting rows affected")
	}

	if rowsAffected == 0 {
		return errors.New("user anime not found")
	}

	return nil
}

func (r *UserAnimeRepository) DeleteByUserAndAnimeMALID(ctx context.Context, userID uint, animeMALID int64) error {
	query := `DELETE FROM user_animes WHERE user_id = $1 AND anime_mal_id = $2`

	result, err := r.db.ExecContext(ctx, query, userID, animeMALID)
	if err != nil {
		r.logger.Error("Error deleting user anime by user ID and anime MAL ID", map[string]interface{}{
			"user_id":      userID,
			"anime_mal_id": animeMALID,
			"error":        err.Error(),
		})
		return errors.Wrap(err, "error deleting user anime by user ID and anime MAL ID")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "error getting rows affected")
	}

	if rowsAffected == 0 {
		return errors.New("user anime not found")
	}

	return nil
}

func (r *UserAnimeRepository) ChangeStatus(ctx context.Context, userID uint, animeMALID int64, status models.WatchStatus) error {
	userAnime, err := r.GetByUserAndAnimeMALID(ctx, userID, animeMALID)
	if err != nil && !strings.Contains(err.Error(), "not found") {
		return err
	}

	now := time.Now()

	if userAnime != nil {
		userAnime.Status = status
		userAnime.UpdatedAt = now
		return r.Update(ctx, userAnime)
	}

	newUserAnime := &models.UserAnime{
		UserID:     userID,
		AnimeMALID: animeMALID,
		Status:     status,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	return r.Create(ctx, newUserAnime)
}

func (r *UserAnimeRepository) UpdateEpisodesWatched(ctx context.Context, userID uint, animeMALID int64, episodesWatched int) error {
	userAnime, err := r.GetByUserAndAnimeMALID(ctx, userID, animeMALID)
	if err != nil {
		return err
	}

	userAnime.EpisodesWatched = episodesWatched
	userAnime.UpdatedAt = time.Now()

	return r.Update(ctx, userAnime)
}

func (r *UserAnimeRepository) UpdateRating(ctx context.Context, userID uint, animeMALID int64, rating float32) error {
	userAnime, err := r.GetByUserAndAnimeMALID(ctx, userID, animeMALID)
	if err != nil {
		return err
	}

	userAnime.Rating = rating
	userAnime.UpdatedAt = time.Now()

	return r.Update(ctx, userAnime)
}

func (r *UserAnimeRepository) GetUserStats(ctx context.Context, userID uint) (*models.AnimeStats, error) {
	query := `
		SELECT 
			COUNT(CASE WHEN status = 'watched' THEN 1 END) as total_watched,
			COUNT(CASE WHEN status = 'plan_to_watch' THEN 1 END) as total_plan_to_watch,
			COUNT(CASE WHEN status = 'watching' THEN 1 END) as total_watching,
			COUNT(CASE WHEN status = 'waiting' THEN 1 END) as total_waiting,
			SUM(episodes_watched) as total_episodes,
			AVG(CASE WHEN rating > 0 THEN rating ELSE NULL END) as average_rating
		FROM user_animes
		WHERE user_id = $1
	`

	stats := &models.AnimeStats{}
	
	var avgRating sql.NullFloat64
	var totalEpisodes sql.NullInt64

	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&stats.TotalWatched,
		&stats.TotalPlanToWatch,
		&stats.TotalWatching,
		&stats.TotalWaiting,
		&totalEpisodes,
		&avgRating,
	)

	if err != nil {
		r.logger.Error("Error getting user anime stats", map[string]interface{}{
			"user_id": userID,
			"error":   err.Error(),
		})
		return nil, errors.Wrap(err, "error getting user anime stats")
	}

	if totalEpisodes.Valid {
		stats.TotalEpisodes = int(totalEpisodes.Int64)
	} else {
		stats.TotalEpisodes = 0
	}

	if avgRating.Valid {
		stats.AverageRating = avgRating.Float64
	} else {
		stats.AverageRating = 0
	}

	return stats, nil
}

func (r *UserAnimeRepository) GetUserAnimeWithDetails(ctx context.Context, filter models.UserAnimeFilter, jikanClient *api.JikanClient) (*models.UserAnimeList, error) {
	userAnimes, total, err := r.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	response := &models.UserAnimeList{
		TotalCount: total,
		Page:       filter.Page,
		Limit:      filter.Limit,
		Items:      make([]*models.UserAnimeWithDetails, 0, len(userAnimes)),
	}

	for _, userAnime := range userAnimes {
		animeDetails, err := jikanClient.GetAnimeByID(ctx, userAnime.AnimeMALID)
		if err != nil {
			r.logger.Warn("Failed to get anime details from Jikan API", map[string]interface{}{
				"anime_mal_id": userAnime.AnimeMALID,
				"error":        err.Error(),
			})
			
			detailedAnime := &models.UserAnimeWithDetails{
				UserAnime:   *userAnime,
				AnimeTitle:  "Unknown",
				AnimeImage:  "",
				AnimeType:   "",
				AnimeStatus: "",
			}
			
			response.Items = append(response.Items, detailedAnime)
			continue
		}

		detailedAnime := &models.UserAnimeWithDetails{
			UserAnime:      *userAnime,
			AnimeTitle:     animeDetails.Title,
			AnimeImage:     animeDetails.ImageURL,
			AnimeType:      animeDetails.Type,
			AnimeEpisodes:  animeDetails.Episodes,
			AnimeStatus:    animeDetails.Status,
			AnimeScore:     animeDetails.Score,
		}

		response.Items = append(response.Items, detailedAnime)
	}

	return response, nil
}

func (r *UserAnimeRepository) CreateOrUpdateUserAnime(ctx context.Context, userAnime *models.UserAnime) error {
	existing, err := r.GetByUserAndAnimeMALID(ctx, userAnime.UserID, userAnime.AnimeMALID)
	if err != nil && !strings.Contains(err.Error(), "not found") {
		return err
	}

	if existing != nil {
		userAnime.ID = existing.ID
		userAnime.CreatedAt = existing.CreatedAt
		userAnime.UpdatedAt = time.Now()
		return r.Update(ctx, userAnime)
	}

	return r.Create(ctx, userAnime)
}