package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/merdernoty/anime-service/internal/domain/models"
	"logur.dev/logur"
)

const (
	jikanBaseURL = "https://api.jikan.moe/v4"
)

type JikanClient struct {
	httpClient *http.Client
	logger     logur.LoggerFacade
}

func NewJikanClient(logger logur.LoggerFacade) *JikanClient {
	return &JikanClient{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		logger: logger,
	}
}

func (c *JikanClient) GetAnimeByID(ctx context.Context, malID int64) (*models.Anime, error) {
	c.logger.Info("Fetching anime from Jikan API", map[string]interface{}{
		"mal_id": malID,
	})

	time.Sleep(500 * time.Millisecond)

	apiURL := fmt.Sprintf("%s/anime/%d", jikanBaseURL, malID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.Error("Error fetching anime from Jikan API", map[string]interface{}{
			"mal_id": malID,
			"error":  err.Error(),
		})
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.logger.Error("Received non-OK response from Jikan API", map[string]interface{}{
			"mal_id":      malID,
			"status_code": resp.StatusCode,
		})
		return nil, fmt.Errorf("received non-OK response: %d", resp.StatusCode)
	}

	var animeResponse struct {
		Data struct {
			MalID         int     `json:"mal_id"`
			Title         string  `json:"title"`
			TitleEnglish  string  `json:"title_english"`
			TitleJapanese string  `json:"title_japanese"`
			Synopsis      string  `json:"synopsis"`
			ImageURL      string  `json:"image_url"`
			Images        struct {
				JPG struct {
					ImageURL string `json:"image_url"`
				} `json:"jpg"`
			} `json:"images"`
			Type       string  `json:"type"`
			Source     string  `json:"source"`
			Episodes   int     `json:"episodes"`
			Status     string  `json:"status"`
			Airing     bool    `json:"airing"`
			Score      float64 `json:"score"`
			Rank       int     `json:"rank"`
			Popularity int     `json:"popularity"`
			Genres     []struct {
				MalID int    `json:"mal_id"`
				Name  string `json:"name"`
			} `json:"genres"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&animeResponse); err != nil {
		c.logger.Error("Failed to decode response", map[string]interface{}{
			"mal_id": malID,
			"error":  err.Error(),
		})
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	anime := &models.Anime{
		MALId:         int64(animeResponse.Data.MalID),
		Title:         animeResponse.Data.Title,
		TitleEnglish:  animeResponse.Data.TitleEnglish,
		TitleJapanese: animeResponse.Data.TitleJapanese,
		Synopsis:      animeResponse.Data.Synopsis,
		ImageURL:      animeResponse.Data.Images.JPG.ImageURL,
		Type:          animeResponse.Data.Type,
		Source:        animeResponse.Data.Source,
		Episodes:      animeResponse.Data.Episodes,
		Status:        animeResponse.Data.Status,
		Airing:        animeResponse.Data.Airing,
		Score:         animeResponse.Data.Score,
		Rank:          animeResponse.Data.Rank,
		Popularity:    animeResponse.Data.Popularity,
	}

	anime.Genres = make([]models.Genre, 0, len(animeResponse.Data.Genres))
	for _, genre := range animeResponse.Data.Genres {
		anime.Genres = append(anime.Genres, models.Genre{
			ID:   int64(genre.MalID),
			Name: genre.Name,
		})
	}

	return anime, nil
}

func (c *JikanClient) SearchAnime(ctx context.Context, query string, page, limit int) ([]*models.Anime, int, error) {
	c.logger.Info("Searching anime in Jikan API", map[string]interface{}{
		"query": query,
		"page":  page,
		"limit": limit,
	})

	time.Sleep(500 * time.Millisecond)

	apiURL, err := url.Parse(fmt.Sprintf("%s/anime", jikanBaseURL))
	if err != nil {
		return nil, 0, fmt.Errorf("failed to parse URL: %w", err)
	}

	q := apiURL.Query()
	q.Add("q", query)
	q.Add("page", strconv.Itoa(page))
	q.Add("limit", strconv.Itoa(limit))
	apiURL.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL.String(), nil)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.Error("Error searching anime from Jikan API", map[string]interface{}{
			"query": query,
			"error": err.Error(),
		})
		return nil, 0, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.logger.Error("Received non-OK response from Jikan API", map[string]interface{}{
			"query":       query,
			"status_code": resp.StatusCode,
		})
		return nil, 0, fmt.Errorf("received non-OK response: %d", resp.StatusCode)
	}

	var searchResponse struct {
		Pagination struct {
			LastVisiblePage int `json:"last_visible_page"`
			HasNextPage     bool `json:"has_next_page"`
		} `json:"pagination"`
		Data []struct {
			MalID    int     `json:"mal_id"`
			Title    string  `json:"title"`
			Synopsis string  `json:"synopsis"`
			Type     string  `json:"type"`
			Episodes int     `json:"episodes"`
			Score    float64 `json:"score"`
			Airing   bool    `json:"airing"`
			Images   struct {
				JPG struct {
					ImageURL string `json:"image_url"`
				} `json:"jpg"`
			} `json:"images"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&searchResponse); err != nil {
		c.logger.Error("Failed to decode response", map[string]interface{}{
			"query": query,
			"error": err.Error(),
		})
		return nil, 0, fmt.Errorf("failed to decode response: %w", err)
	}

	animes := make([]*models.Anime, 0, len(searchResponse.Data))
	for _, result := range searchResponse.Data {
		anime := &models.Anime{
			MALId:     int64(result.MalID),
			Title:     result.Title,
			Synopsis:  result.Synopsis,
			Type:      result.Type,
			Episodes:  result.Episodes,
			Score:     result.Score,
			Airing:    result.Airing,
			ImageURL:  result.Images.JPG.ImageURL,
		}
		animes = append(animes, anime)
	}

	return animes, searchResponse.Pagination.LastVisiblePage, nil
}

func (c *JikanClient) GetTopAnime(ctx context.Context, page, limit int) ([]*models.Anime, int, error) {
	c.logger.Info("Fetching top anime from Jikan API", map[string]interface{}{
		"page":  page,
		"limit": limit,
	})

	time.Sleep(500 * time.Millisecond)

	apiURL, err := url.Parse(fmt.Sprintf("%s/top/anime", jikanBaseURL))
	if err != nil {
		return nil, 0, fmt.Errorf("failed to parse URL: %w", err)
	}

	q := apiURL.Query()
	q.Add("page", strconv.Itoa(page))
	q.Add("limit", strconv.Itoa(limit))
	apiURL.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL.String(), nil)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.Error("Error fetching top anime from Jikan API", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, 0, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.logger.Error("Received non-OK response from Jikan API", map[string]interface{}{
			"status_code": resp.StatusCode,
		})
		return nil, 0, fmt.Errorf("received non-OK response: %d", resp.StatusCode)
	}

	var topResponse struct {
		Pagination struct {
			LastVisiblePage int `json:"last_visible_page"`
			HasNextPage     bool `json:"has_next_page"`
		} `json:"pagination"`
		Data []struct {
			MalID    int     `json:"mal_id"`
			Title    string  `json:"title"`
			Type     string  `json:"type"`
			Episodes int     `json:"episodes"`
			Score    float64 `json:"score"`
			Rank     int     `json:"rank"`
			Images   struct {
				JPG struct {
					ImageURL string `json:"image_url"`
				} `json:"jpg"`
			} `json:"images"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&topResponse); err != nil {
		c.logger.Error("Failed to decode response", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, 0, fmt.Errorf("failed to decode response: %w", err)
	}

	animes := make([]*models.Anime, 0, len(topResponse.Data))
	for _, result := range topResponse.Data {
		anime := &models.Anime{
			MALId:    int64(result.MalID),
			Title:    result.Title,
			Type:     result.Type,
			Episodes: result.Episodes,
			Score:    result.Score,
			Rank:     result.Rank,
			ImageURL: result.Images.JPG.ImageURL,
		}
		animes = append(animes, anime)
	}

	return animes, topResponse.Pagination.LastVisiblePage, nil
}

func (c *JikanClient) GetSeasonalAnime(ctx context.Context, year, season string, page, limit int) ([]*models.Anime, int, error) {
	c.logger.Info("Fetching seasonal anime from Jikan API", map[string]interface{}{
		"year":   year,
		"season": season,
		"page":   page,
		"limit":  limit,
	})

	time.Sleep(500 * time.Millisecond)

	apiURL, err := url.Parse(fmt.Sprintf("%s/seasons/%s/%s", jikanBaseURL, year, season))
	if err != nil {
		return nil, 0, fmt.Errorf("failed to parse URL: %w", err)
	}

	q := apiURL.Query()
	q.Add("page", strconv.Itoa(page))
	q.Add("limit", strconv.Itoa(limit))
	apiURL.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL.String(), nil)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.Error("Error fetching seasonal anime from Jikan API", map[string]interface{}{
			"year":   year,
			"season": season,
			"error":  err.Error(),
		})
		return nil, 0, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.logger.Error("Received non-OK response from Jikan API", map[string]interface{}{
			"year":        year,
			"season":      season,
			"status_code": resp.StatusCode,
		})
		return nil, 0, fmt.Errorf("received non-OK response: %d", resp.StatusCode)
	}

	var seasonalResponse struct {
		Pagination struct {
			LastVisiblePage int `json:"last_visible_page"`
			HasNextPage     bool `json:"has_next_page"`
		} `json:"pagination"`
		Data []struct {
			MalID    int     `json:"mal_id"`
			Title    string  `json:"title"`
			Synopsis string  `json:"synopsis"`
			Type     string  `json:"type"`
			Episodes int     `json:"episodes"`
			Score    float64 `json:"score"`
			Airing   bool    `json:"airing"`
			Images   struct {
				JPG struct {
					ImageURL string `json:"image_url"`
				} `json:"jpg"`
			} `json:"images"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&seasonalResponse); err != nil {
		c.logger.Error("Failed to decode response", map[string]interface{}{
			"year":   year,
			"season": season,
			"error":  err.Error(),
		})
		return nil, 0, fmt.Errorf("failed to decode response: %w", err)
	}

	animes := make([]*models.Anime, 0, len(seasonalResponse.Data))
	for _, result := range seasonalResponse.Data {
		anime := &models.Anime{
			MALId:     int64(result.MalID),
			Title:     result.Title,
			Synopsis:  result.Synopsis,
			Type:      result.Type,
			Episodes:  result.Episodes,
			Score:     result.Score,
			Airing:    result.Airing,
			ImageURL:  result.Images.JPG.ImageURL,
		}
		animes = append(animes, anime)
	}

	return animes, seasonalResponse.Pagination.LastVisiblePage, nil
}

func (c *JikanClient) GetAnimeRecommendations(ctx context.Context, malID int64, page, limit int) ([]*models.Anime, error) {
	c.logger.Info("Fetching anime recommendations from Jikan API", map[string]interface{}{
		"mal_id": malID,
	})

	time.Sleep(500 * time.Millisecond)

	apiURL := fmt.Sprintf("%s/anime/%d/recommendations", jikanBaseURL, malID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.Error("Error fetching anime recommendations from Jikan API", map[string]interface{}{
			"mal_id": malID,
			"error":  err.Error(),
		})
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.logger.Error("Received non-OK response from Jikan API", map[string]interface{}{
			"mal_id":      malID,
			"status_code": resp.StatusCode,
		})
		return nil, fmt.Errorf("received non-OK response: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.logger.Error("Failed to read response body", map[string]interface{}{
			"mal_id": malID,
			"error":  err.Error(),
		})
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var recommendationsResponse struct {
		Data []struct {
			Entry struct {
				MalID  int    `json:"mal_id"`
				Title  string `json:"title"`
				Images struct {
					JPG struct {
						ImageURL string `json:"image_url"`
					} `json:"jpg"`
				} `json:"images"`
			} `json:"entry"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &recommendationsResponse); err != nil {
		c.logger.Error("Failed to decode response", map[string]interface{}{
			"mal_id": malID,
			"error":  err.Error(),
		})
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	allRecommendations := recommendationsResponse.Data
	startIdx := (page - 1) * limit
	endIdx := startIdx + limit

	if startIdx >= len(allRecommendations) {
		return []*models.Anime{}, nil
	}

	if endIdx > len(allRecommendations) {
		endIdx = len(allRecommendations)
	}

	pagedRecommendations := allRecommendations[startIdx:endIdx]

	animes := make([]*models.Anime, 0, len(pagedRecommendations))
	for _, rec := range pagedRecommendations {
		anime := &models.Anime{
			MALId:    int64(rec.Entry.MalID),
			Title:    rec.Entry.Title,
			ImageURL: rec.Entry.Images.JPG.ImageURL,
		}
		animes = append(animes, anime)
	}

	return animes, nil
}