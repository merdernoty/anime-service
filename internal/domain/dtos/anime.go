package dtos
import "github.com/merdernoty/anime-service/internal/domain/models"

type AnimeResponse struct {
	MALId         int64         `json:"mal_id" example:"5114"`
	Title         string        `json:"title" example:"Fullmetal Alchemist: Brotherhood"`
	TitleEnglish  string        `json:"title_english,omitempty" example:"Fullmetal Alchemist: Brotherhood"`
	TitleJapanese string        `json:"title_japanese,omitempty" example:"鋼の錬金術師 FULLMETAL ALCHEMIST"`
	Synopsis      string        `json:"synopsis,omitempty" example:"After a terrible alchemy experiment gone wrong..."`
	ImageURL      string        `json:"image_url" example:"https://cdn.myanimelist.net/images/anime/1223/96541.jpg"`
	Type          string        `json:"type,omitempty" example:"TV"`
	Source        string        `json:"source,omitempty" example:"Manga"`
	Episodes      int           `json:"episodes,omitempty" example:"64"`
	Status        string        `json:"status,omitempty" example:"Finished Airing"`
	Airing        bool          `json:"airing" example:"false"`
	Score         float64       `json:"score,omitempty" example:"9.16"`
	Rank          int           `json:"rank,omitempty" example:"1"`
	Popularity    int           `json:"popularity,omitempty" example:"3"`
	Genres        []GenreObject `json:"genres,omitempty"`
}

type GenreObject struct {
	ID   int64  `json:"id" example:"1"`
	Name string `json:"name" example:"Action"`
}

type AnimeListResponse struct {
	Items      []AnimeResponse `json:"items"`
	Page       int             `json:"page" example:"1"`
	Limit      int             `json:"limit" example:"10"`
	TotalPages int             `json:"total_pages" example:"50"`
}

type UserAnimeResponse struct {
	ID             uint              `json:"id" example:"1"`
	UserID         uint              `json:"user_id" example:"42"`
	AnimeMALID     int64             `json:"anime_mal_id" example:"5114"`
	Status         models.WatchStatus `json:"status" example:"watching"`
	Rating         float32           `json:"rating" example:"9.5"`
	Notes          string            `json:"notes,omitempty" example:"My favorite anime!"`
	EpisodesWatched int              `json:"episodes_watched" example:"24"`
	AnimeTitle     string            `json:"anime_title" example:"Fullmetal Alchemist: Brotherhood"`
	AnimeImage     string            `json:"anime_image" example:"https://cdn.myanimelist.net/images/anime/1223/96541.jpg"`
	AnimeType      string            `json:"anime_type" example:"TV"`
	AnimeEpisodes  int               `json:"anime_episodes" example:"64"`
	AnimeStatus    string            `json:"anime_status" example:"Finished Airing"`
	AnimeScore     float64           `json:"anime_score" example:"9.16"`
}

type UserAnimeListResponse struct {
	Items      []UserAnimeResponse `json:"items"`
	TotalCount int                 `json:"total_count" example:"42"`
	Page       int                 `json:"page" example:"1"`
	Limit      int                 `json:"limit" example:"10"`
}

type StatsResponse struct {
	TotalWatched     int     `json:"total_watched" example:"15"`
	TotalPlanToWatch int     `json:"total_plan_to_watch" example:"30"`
	TotalWatching    int     `json:"total_watching" example:"5"`
	TotalWaiting     int     `json:"total_waiting" example:"10"`
	TotalEpisodes    int     `json:"total_episodes" example:"347"`
	AverageRating    float64 `json:"average_rating" example:"8.75"`
}

type AddAnimeRequest struct {
	AnimeMALID int64              `json:"anime_mal_id" binding:"required" example:"5114"`
	Status     models.WatchStatus `json:"status" binding:"required" example:"watching"`
}

type UpdateStatusRequest struct {
	Status models.WatchStatus `json:"status" binding:"required" example:"watched"`
}

type UpdateEpisodesRequest struct {
	EpisodesWatched int `json:"episodes_watched" binding:"required" example:"24"`
}

type UpdateRatingRequest struct {
	Rating float32 `json:"rating" binding:"required,min=0,max=10" example:"9.5"`
}
