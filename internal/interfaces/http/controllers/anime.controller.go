package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/merdernoty/anime-service/internal/application/services"
	"github.com/merdernoty/anime-service/internal/domain/dtos"
	"github.com/merdernoty/anime-service/internal/domain/models"
	"logur.dev/logur"
)

type AnimeController struct {
	animeService services.AnimeServiceImpl
	logger       logur.LoggerFacade
}

func NewAnimeController(animeService services.AnimeServiceImpl, logger logur.LoggerFacade) *AnimeController {
	return &AnimeController{
		animeService: animeService,
		logger:       logger,
	}
}

func handleAnimeError(ctx *gin.Context, err error) {
	switch {
	case err == services.ErrAnimeAlreadyExists:
		ctx.JSON(http.StatusConflict, gin.H{
			"error":   "anime already exists",
			"details": err.Error(),
		})
	case err == services.ErrAnimeDeleteFailed:
		ctx.JSON(http.StatusConflict, gin.H{
			"error":   "anime delete failed",
			"details": err.Error(),
		})
	case err == services.ErrAnimeNotInUserList:
		ctx.JSON(http.StatusNotFound, gin.H{
			"error":   "anime not found in user list",
			"details": err.Error(),
		})
	case err == services.ErrAnimeUpdateFailed:
		ctx.JSON(http.StatusConflict, gin.H{
			"error":   "anime update failed",
			"details": err.Error(),
		})
	case err == services.ErrAnimeStatsFailed:
		ctx.JSON(http.StatusNotFound, gin.H{
			"error":   "anime stats not found",
			"details": err.Error(),
		})
	case err == services.ErrAnimeNotFound:
		ctx.JSON(http.StatusNotFound, gin.H{
			"error":   "anime not found",
			"details": err.Error(),
		})
	case err == services.ErrFetchAnimeFailed:
		ctx.JSON(http.StatusNotFound, gin.H{
			"error":   "fetch anime failed",
			"details": err.Error(),
		})
	default:
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "internal server error",
			"details": err.Error(),
		})
	}
}




// GetAnimeByID godoc
//	@Summary		Получить информацию об аниме по его ID
//	@Description	Получает детальную информацию об аниме по его MAL ID
//	@Tags			anime
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"MAL ID аниме"
//	@Success		200	{object}	dtos.AnimeResponse
//	@Failure		400	{object}	map[string]string	"Неверный ID аниме"
//	@Failure		404	{object}	map[string]string	"Аниме не найдено"
//	@Failure		500	{object}	map[string]string	"Внутренняя ошибка сервера"
//	@Router			/anime/{id} [get]
func (c *AnimeController) GetAnimeByID(ctx *gin.Context) {
	malID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID аниме"})
		return
	}

	anime, err := c.animeService.GetAnimeByID(ctx, malID)
	if err != nil {
		c.logger.Error("Error getting anime by ID", map[string]interface{}{
			"mal_id": malID,
			"error":  err.Error(),
		})
		
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось получить информацию об аниме"})
		return
	}

	if anime == nil {
		handleAnimeError(ctx, services.ErrAnimeNotFound)
		return
	}

	// Преобразуем модель аниме в ответ API
	genres := make([]dtos.GenreObject, 0, len(anime.Genres))
	for _, genre := range anime.Genres {
		genres = append(genres, dtos.GenreObject{
			ID:   genre.ID,
			Name: genre.Name,
		})
	}

	response := dtos.AnimeResponse{
		MALId:         anime.MALId,
		Title:         anime.Title,
		TitleEnglish:  anime.TitleEnglish,
		TitleJapanese: anime.TitleJapanese,
		Synopsis:      anime.Synopsis,
		ImageURL:      anime.ImageURL,
		Type:          anime.Type,
		Source:        anime.Source,
		Episodes:      anime.Episodes,
		Status:        anime.Status,
		Airing:        anime.Airing,
		Score:         anime.Score,
		Rank:          anime.Rank,
		Popularity:    anime.Popularity,
		Genres:        genres,
	}

	ctx.JSON(http.StatusOK, response)
}

// SearchAnime godoc
//	@Summary		Поиск аниме по запросу
//	@Description	Выполняет поиск аниме по заданному запросу с пагинацией
//	@Tags			anime
//	@Accept			json
//	@Produce		json
//	@Param			query	query		string	true	"Поисковый запрос"
//	@Param			page	query		int		false	"Номер страницы"						default(1)	minimum(1)
//	@Param			limit	query		int		false	"Количество результатов на странице"	default(10)	minimum(1)	maximum(50)
//	@Success		200		{object}	dtos.AnimeListResponse
//	@Failure		500		{object}	map[string]string	"Внутренняя ошибка сервера"
//	@Router			/anime/search [get]
func (c *AnimeController) SearchAnime(ctx *gin.Context) {
	query := ctx.Query("query")
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	// Проверяем валидность параметров
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	} else if limit > 50 {
		limit = 50
	}

	animes, totalPages, err := c.animeService.SearchAnime(ctx, query, page, limit)
	if err != nil {
		c.logger.Error("Error searching anime", map[string]interface{}{
			"query": query,
			"error": err.Error(),
		})
		handleAnimeError(ctx, err)
		return
	}

	// Преобразуем модели аниме в ответ API
	items := make([]dtos.AnimeResponse, 0, len(animes))
	for _, anime := range animes {
		genres := make([]dtos.GenreObject, 0, len(anime.Genres))
		for _, genre := range anime.Genres {
			genres = append(genres, dtos.GenreObject{
				ID:   genre.ID,
				Name: genre.Name,
			})
		}

		items = append(items, dtos.AnimeResponse{
			MALId:         anime.MALId,
			Title:         anime.Title,
			TitleEnglish:  anime.TitleEnglish,
			TitleJapanese: anime.TitleJapanese,
			Synopsis:      anime.Synopsis,
			ImageURL:      anime.ImageURL,
			Type:          anime.Type,
			Source:        anime.Source,
			Episodes:      anime.Episodes,
			Status:        anime.Status,
			Airing:        anime.Airing,
			Score:         anime.Score,
			Rank:          anime.Rank,
			Popularity:    anime.Popularity,
			Genres:        genres,
		})
	}

	ctx.JSON(http.StatusOK, dtos.AnimeListResponse{
		Items:      items,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	})
}

// GetTopAnime godoc
//	@Summary		Получить список популярных аниме
//	@Description	Возвращает список популярных аниме с пагинацией
//	@Tags			anime
//	@Accept			json
//	@Produce		json
//	@Param			page	query		int	false	"Номер страницы"						default(1)	minimum(1)
//	@Param			limit	query		int	false	"Количество результатов на странице"	default(10)	minimum(1)	maximum(50)
//	@Success		200		{object}	dtos.AnimeListResponse
//	@Failure		500		{object}	map[string]string	"Внутренняя ошибка сервера"
//	@Router			/anime/top [get]
func (c *AnimeController) GetTopAnime(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	// Проверяем валидность параметров
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	} else if limit > 50 {
		limit = 50
	}

	animes, totalPages, err := c.animeService.GetTopAnime(ctx, page, limit)
	if err != nil {
		c.logger.Error("Error getting top anime", map[string]interface{}{
			"error": err.Error(),
		})
		handleAuthError(ctx, err)
		return
	}

	// Преобразуем модели аниме в ответ API
	items := make([]dtos.AnimeResponse, 0, len(animes))
	for _, anime := range animes {
		genres := make([]dtos.GenreObject, 0, len(anime.Genres))
		for _, genre := range anime.Genres {
			genres = append(genres, dtos.GenreObject{
				ID:   genre.ID,
				Name: genre.Name,
			})
		}

		items = append(items, dtos.AnimeResponse{
			MALId:         anime.MALId,
			Title:         anime.Title,
			TitleEnglish:  anime.TitleEnglish,
			TitleJapanese: anime.TitleJapanese,
			Synopsis:      anime.Synopsis,
			ImageURL:      anime.ImageURL,
			Type:          anime.Type,
			Source:        anime.Source,
			Episodes:      anime.Episodes,
			Status:        anime.Status,
			Airing:        anime.Airing,
			Score:         anime.Score,
			Rank:          anime.Rank,
			Popularity:    anime.Popularity,
			Genres:        genres,
		})
	}

	ctx.JSON(http.StatusOK, dtos.AnimeListResponse{
		Items:      items,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	})
}

// GetSeasonalAnime godoc
//	@Summary		Получить список сезонных аниме
//	@Description	Возвращает список аниме для указанного сезона и года
//	@Tags			anime
//	@Accept			json
//	@Produce		json
//	@Param			year	path		string	true	"Год (например, 2023)"
//	@Param			season	path		string	true	"Сезон (winter, spring, summer, fall)"
//	@Param			page	query		int		false	"Номер страницы"						default(1)	minimum(1)
//	@Param			limit	query		int		false	"Количество результатов на странице"	default(10)	minimum(1)	maximum(50)
//	@Success		200		{object}	dtos.AnimeListResponse
//	@Failure		400		{object}	map[string]string	"Неверные параметры запроса"
//	@Failure		500		{object}	map[string]string	"Внутренняя ошибка сервера"
//	@Router			/anime/seasonal/{year}/{season} [get]
func (c *AnimeController) GetSeasonalAnime(ctx *gin.Context) {
	year := ctx.Param("year")
	season := ctx.Param("season")
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	// Проверяем валидность параметров
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	} else if limit > 50 {
		limit = 50
	}

	// Проверяем валидность сезона
	validSeasons := map[string]bool{
		"winter": true,
		"spring": true,
		"summer": true,
		"fall":   true,
	}
	if !validSeasons[season] {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Неверный сезон. Допустимые значения: winter, spring, summer, fall"})
		return
	}

	animes, totalPages, err := c.animeService.GetSeasonalAnime(ctx, year, season, page, limit)
	if err != nil {
		c.logger.Error("Error getting seasonal anime", map[string]interface{}{
			"year":   year,
			"season": season,
			"error":  err.Error(),
		})
		handleAuthError(ctx, err)
		return
	}

	// Преобразуем модели аниме в ответ API
	items := make([]dtos.AnimeResponse, 0, len(animes))
	for _, anime := range animes {
		genres := make([]dtos.GenreObject, 0, len(anime.Genres))
		for _, genre := range anime.Genres {
			genres = append(genres, dtos.GenreObject{
				ID:   genre.ID,
				Name: genre.Name,
			})
		}

		items = append(items, dtos.AnimeResponse{
			MALId:         anime.MALId,
			Title:         anime.Title,
			TitleEnglish:  anime.TitleEnglish,
			TitleJapanese: anime.TitleJapanese,
			Synopsis:      anime.Synopsis,
			ImageURL:      anime.ImageURL,
			Type:          anime.Type,
			Source:        anime.Source,
			Episodes:      anime.Episodes,
			Status:        anime.Status,
			Airing:        anime.Airing,
			Score:         anime.Score,
			Rank:          anime.Rank,
			Popularity:    anime.Popularity,
			Genres:        genres,
		})
	}

	ctx.JSON(http.StatusOK, dtos.AnimeListResponse{
		Items:      items,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	})
}

// GetAnimeRecommendations godoc
//	@Summary		Получить рекомендации аниме
//	@Description	Возвращает список рекомендаций аниме на основе указанного MAL ID
//	@Tags			anime
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int	true	"MAL ID аниме для получения рекомендаций"
//	@Param			page	query		int	false	"Номер страницы"						default(1)	minimum(1)
//	@Param			limit	query		int	false	"Количество результатов на странице"	default(10)	minimum(1)	maximum(50)
//	@Success		200		{object}	dtos.AnimeListResponse
//	@Failure		400		{object}	map[string]string	"Неверный ID аниме"
//	@Failure		500		{object}	map[string]string	"Внутренняя ошибка сервера"
//	@Router			/anime/{id}/recommendations [get]
func (c *AnimeController) GetAnimeRecommendations(ctx *gin.Context) {
	malID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID аниме"})
		return
	}

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	// Проверяем валидность параметров
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	} else if limit > 50 {
		limit = 50
	}

	animes, err := c.animeService.GetAnimeRecommendations(ctx, malID, page, limit)
	if err != nil {
		c.logger.Error("Error getting anime recommendations", map[string]interface{}{
			"mal_id": malID,
			"error":  err.Error(),
		})
		handleAuthError(ctx, err)
		return
	}

	// Преобразуем модели аниме в ответ API
	items := make([]dtos.AnimeResponse, 0, len(animes))
	for _, anime := range animes {
		items = append(items, dtos.AnimeResponse{
			MALId:    anime.MALId,
			Title:    anime.Title,
			ImageURL: anime.ImageURL,
		})
	}

	ctx.JSON(http.StatusOK, dtos.AnimeListResponse{
		Items:      items,
		Page:       page,
		Limit:      limit,
		TotalPages: 1, // API не предоставляет информацию о общем количестве страниц
	})
}

// GetUserAnimeList godoc
//	@Summary		Получить список аниме пользователя
//	@Description	Возвращает список аниме пользователя с возможностью фильтрации по статусу
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			user_id	path		int		true	"ID пользователя"
//	@Param			status	query		string	false	"Статус аниме (watched, plan_to_watch, watching, waiting)"
//	@Param			page	query		int		false	"Номер страницы"						default(1)	minimum(1)
//	@Param			limit	query		int		false	"Количество результатов на странице"	default(10)	minimum(1)	maximum(50)
//	@Success		200		{object}	dtos.UserAnimeListResponse
//	@Failure		400		{object}	map[string]string	"Неверный ID пользователя"
//	@Failure		500		{object}	map[string]string	"Внутренняя ошибка сервера"
//	@Router			/users/{user_id}/anime [get]
func (c *AnimeController) GetUserAnimeList(ctx *gin.Context) {
	userIDStr := ctx.Param("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID пользователя"})
		return
	}

	status := models.WatchStatus(ctx.DefaultQuery("status", ""))
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	query := ctx.DefaultQuery("query", "")

	// Проверяем валидность параметров
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	} else if limit > 50 {
		limit = 50
	}

	filter := models.UserAnimeFilter{
		UserID: int64(userID),
		Status: status,
		Page:   page,
		Limit:  limit,
		Query:  query,
	}

	userAnimeList, err := c.animeService.GetUserAnimeList(ctx, filter)
	if err != nil {
		c.logger.Error("Error getting user anime list", map[string]interface{}{
			"user_id": userID,
			"error":   err.Error(),
		})
		handleAuthError(ctx, err)
		return
	}

	// Преобразуем результаты в ответ API
	items := make([]dtos.UserAnimeResponse, 0, len(userAnimeList.Items))
	for _, item := range userAnimeList.Items {
		items = append(items, dtos.UserAnimeResponse{
			ID:              item.ID,
			UserID:          item.UserID,
			AnimeMALID:      item.AnimeMALID,
			Status:          item.Status,
			Rating:          item.Rating,
			Notes:           item.Notes,
			EpisodesWatched: item.EpisodesWatched,
			AnimeTitle:      item.AnimeTitle,
			AnimeImage:      item.AnimeImage,
			AnimeType:       item.AnimeType,
			AnimeEpisodes:   item.AnimeEpisodes,
			AnimeStatus:     item.AnimeStatus,
			AnimeScore:      item.AnimeScore,
		})
	}

	ctx.JSON(http.StatusOK, dtos.UserAnimeListResponse{
		Items:      items,
		TotalCount: userAnimeList.TotalCount,
		Page:       userAnimeList.Page,
		Limit:      userAnimeList.Limit,
	})
}

// AddAnimeToUserList godoc
//	@Summary		Добавить аниме в список пользователя
//	@Description	Добавляет аниме в список пользователя с указанным статусом
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			user_id	path		int						true	"ID пользователя"
//	@Param			anime	body		dtos.AddAnimeRequest	true	"Информация о добавляемом аниме"
//	@Success		201		{object}	map[string]string		"Успешное добавление"
//	@Failure		400		{object}	map[string]string		"Неверные входные данные"
//	@Failure		404		{object}	map[string]string		"Аниме не найдено"
//	@Failure		500		{object}	map[string]string		"Внутренняя ошибка сервера"
//	@Router			/users/{user_id}/anime [post]
func (c *AnimeController) AddAnimeToUserList(ctx *gin.Context) {
	userIDStr := ctx.Param("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		handleAuthError(ctx, err)
		return
	}

	var request dtos.AddAnimeRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		handleAuthError(ctx, err)
		return
	}

	err = c.animeService.AddAnimeToUserList(ctx, uint(userID), request.AnimeMALID, request.Status)
	if err != nil {
		c.logger.Error("Error adding anime to user list", map[string]interface{}{
			"user_id":      userID,
			"anime_mal_id": request.AnimeMALID,
			"error":        err.Error(),
		})
		handleAuthError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Аниме успешно добавлено в список пользователя"})
}

// RemoveAnimeFromUserList godoc
//	@Summary		Удалить аниме из списка пользователя
//	@Description	Удаляет аниме из списка пользователя
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			user_id		path		int					true	"ID пользователя"
//	@Param			anime_id	path		int					true	"MAL ID аниме"
//	@Success		200			{object}	map[string]string	"Успешное удаление"
//	@Failure		400			{object}	map[string]string	"Неверные входные данные"
//	@Failure		404			{object}	map[string]string	"Запись не найдена"
//	@Failure		500			{object}	map[string]string	"Внутренняя ошибка сервера"
//	@Router			/users/{user_id}/anime/{anime_id} [delete]
func (c *AnimeController) RemoveAnimeFromUserList(ctx *gin.Context) {
	userIDStr := ctx.Param("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		handleAuthError(ctx, err)
		return
	}

	animeIDStr := ctx.Param("anime_id")
	animeMALID, err := strconv.ParseInt(animeIDStr, 10, 64)
	if err != nil {
		handleAuthError(ctx, err)
		return
	}

	err = c.animeService.RemoveAnimeFromUserList(ctx, uint(userID), animeMALID)
	if err != nil {
		c.logger.Error("Error removing anime from user list", map[string]interface{}{
			"user_id":      userID,
			"anime_mal_id": animeMALID,
			"error":        err.Error(),
		})
		
		if err.Error() == "user anime not found" {
			handleAuthError(ctx, err)
			return
		}
		
		handleAuthError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Аниме успешно удалено из списка пользователя"})
}

// UpdateUserAnimeStatus godoc
//	@Summary		Обновить статус аниме в списке пользователя
//	@Description	Обновляет статус аниме в списке пользователя
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			user_id		path		int							true	"ID пользователя"
//	@Param			anime_id	path		int							true	"MAL ID аниме"
//	@Param			status		body		dtos.UpdateStatusRequest	true	"Новый статус аниме"
//	@Success		200			{object}	map[string]string			"Успешное обновление"
//	@Failure		400			{object}	map[string]string			"Неверные входные данные"
//	@Failure		404			{object}	map[string]string			"Запись не найдена"
//	@Failure		500			{object}	map[string]string			"Внутренняя ошибка сервера"
//	@Router			/users/{user_id}/anime/{anime_id}/status [put]
func (c *AnimeController) UpdateUserAnimeStatus(ctx *gin.Context) {
	userIDStr := ctx.Param("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		handleAuthError(ctx, err)
		return
	}

	animeIDStr := ctx.Param("anime_id")
	animeMALID, err := strconv.ParseInt(animeIDStr, 10, 64)
	if err != nil {
		handleAuthError(ctx, err)
		return
	}

	var request dtos.UpdateStatusRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		handleAuthError(ctx, err)
		return
	}

	err = c.animeService.UpdateUserAnimeStatus(ctx, uint(userID), animeMALID, request.Status)
	if err != nil {
		c.logger.Error("Error updating user anime status", map[string]interface{}{
			"user_id":      userID,
			"anime_mal_id": animeMALID,
			"error":        err.Error(),
		})
		
		if err.Error() == "user anime not found" {
			handleAuthError(ctx, err)
			return
		}
		
		handleAuthError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Статус аниме успешно обновлен"})
}

// UpdateUserAnimeEpisodes godoc
//	@Summary		Обновить количество просмотренных эпизодов
//	@Description	Обновляет количество просмотренных эпизодов аниме в списке пользователя
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			user_id		path		int							true	"ID пользователя"
//	@Param			anime_id	path		int							true	"MAL ID аниме"
//	@Param			episodes	body		dtos.UpdateEpisodesRequest	true	"Новое количество просмотренных эпизодов"
//	@Success		200			{object}	map[string]string			"Успешное обновление"
//	@Failure		400			{object}	map[string]string			"Неверные входные данные"
//	@Failure		404			{object}	map[string]string			"Запись не найдена"
//	@Failure		500			{object}	map[string]string			"Внутренняя ошибка сервера"
//	@Router			/users/{user_id}/anime/{anime_id}/episodes [put]
func (c *AnimeController) UpdateUserAnimeEpisodes(ctx *gin.Context) {
	userIDStr := ctx.Param("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		handleAuthError(ctx, err)
		return
	}

	animeIDStr := ctx.Param("anime_id")
	animeMALID, err := strconv.ParseInt(animeIDStr, 10, 64)
	if err != nil {
		handleAuthError(ctx, err)
		return
	}

	var request dtos.UpdateEpisodesRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		handleAuthError(ctx, err)
		return
	}

	err = c.animeService.UpdateUserAnimeEpisodes(ctx, uint(userID), animeMALID, request.EpisodesWatched)
	if err != nil {
		c.logger.Error("Error updating user anime episodes", map[string]interface{}{
			"user_id":          userID,
			"anime_mal_id":     animeMALID,
			"episodes_watched": request.EpisodesWatched,
			"error":            err.Error(),
		})
		
		if err.Error() == "user anime not found" {
			handleAuthError(ctx, err)
			return
		}
		
		handleAuthError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Количество просмотренных эпизодов успешно обновлено"})
}

// UpdateUserAnimeRating godoc
//	@Summary		Обновить рейтинг аниме
//	@Description	Обновляет пользовательский рейтинг аниме в списке пользователя
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			user_id		path		int							true	"ID пользователя"
//	@Param			anime_id	path		int							true	"MAL ID аниме"
//	@Param			rating		body		dtos.UpdateRatingRequest	true	"Новый рейтинг аниме (от 0 до 10)"
//	@Success		200			{object}	map[string]string			"Успешное обновление"
//	@Failure		400			{object}	map[string]string			"Неверные входные данные"
//	@Failure		404			{object}	map[string]string			"Запись не найдена"
//	@Failure		500			{object}	map[string]string			"Внутренняя ошибка сервера"
//	@Router			/users/{user_id}/anime/{anime_id}/rating [put]
func (c *AnimeController) UpdateUserAnimeRating(ctx *gin.Context) {
	userIDStr := ctx.Param("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		handleAuthError(ctx, err)
		return
	}

	animeIDStr := ctx.Param("anime_id")
	animeMALID, err := strconv.ParseInt(animeIDStr, 10, 64)
	if err != nil {
		handleAuthError(ctx, err)
		return
	}

	var request dtos.UpdateRatingRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		handleAuthError(ctx, err)
		return
	}

	// Проверяем диапазон рейтинга
	if request.Rating < 0 || request.Rating > 10 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Рейтинг должен быть от 0 до 10"})
		return
	}

	err = c.animeService.UpdateUserAnimeRating(ctx, uint(userID), animeMALID, request.Rating)
	if err != nil {
		c.logger.Error("Error updating user anime rating", map[string]interface{}{
			"user_id":      userID,
			"anime_mal_id": animeMALID,
			"rating":       request.Rating,
			"error":        err.Error(),
		})
		
		if err.Error() == "user anime not found" {
			handleAuthError(ctx, err)
			return
		}
		
		handleAuthError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Рейтинг аниме успешно обновлен"})
}

// GetUserAnimeStats godoc
//	@Summary		Получить статистику пользователя по аниме
//	@Description	Возвращает статистику пользователя по просмотру аниме
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			user_id	path		int					true	"ID пользователя"
//	@Success		200		{object}	dtos.StatsResponse	"Статистика пользователя"
//	@Failure		400		{object}	map[string]string	"Неверный ID пользователя"
//	@Failure		500		{object}	map[string]string	"Внутренняя ошибка сервера"
//	@Router			/users/{user_id}/anime/stats [get]
func (c *AnimeController) GetUserAnimeStats(ctx *gin.Context) {
	userIDStr := ctx.Param("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		handleAuthError(ctx, err)
		return
	}

	stats, err := c.animeService.GetUserAnimeStats(ctx, uint(userID))
	if err != nil {
		c.logger.Error("Error getting user anime stats", map[string]interface{}{
			"user_id": userID,
			"error":   err.Error(),
		})
		handleAuthError(ctx, err)
		return
	}

	response := dtos.StatsResponse{
		TotalWatched:     stats.TotalWatched,
		TotalPlanToWatch: stats.TotalPlanToWatch,
		TotalWatching:    stats.TotalWatching,
		TotalWaiting:     stats.TotalWaiting,
		TotalEpisodes:    stats.TotalEpisodes,
		AverageRating:    stats.AverageRating,
	}

	ctx.JSON(http.StatusOK, response)
}