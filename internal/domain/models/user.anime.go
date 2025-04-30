package models

import (
	"time"
)

type WatchStatus string

const (
	StatusWatched     WatchStatus = "watched"
	StatusPlanToWatch WatchStatus = "plan_to_watch" 
	StatusWatching    WatchStatus = "watching" 
	StatusWaiting     WatchStatus = "waiting"  
)

type UserAnime struct {
	ID		uint       `json:"id" db:"id" gorm:"primaryKey"`
	UserID	uint       `json:"user_id" db:"user_id" gorm:"not null"`
	AnimeMALID  int64       `json:"anime_mal_id" db:"anime_mal_id"`
	Status	WatchStatus `json:"status" db:"status" gorm:"not null"`
	Rating	float32    `json:"rating" db:"rating"`
	Notes       string      `json:"notes" db:"notes"`
	EpisodesWatched int     `json:"episodes_watched" db:"episodes_watched"` 
	CreatedAt   time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at" db:"updated_at"`
}

type UserAnimeFilter struct {
	UserID int64       `json:"user_id" form:"user_id"`
	Status WatchStatus `json:"status" form:"status"`
	Query  string      `json:"query" form:"query"`
	Page   int         `json:"page" form:"page"` 
	Limit  int         `json:"limit" form:"limit"`
}

type UserAnimeList struct {
	Items      []*UserAnimeWithDetails `json:"items"`
	TotalCount int                     `json:"total_count"`
	Page       int                     `json:"page"`
	Limit      int                     `json:"limit"`
}

type UserAnimeWithDetails struct {
	UserAnime    `json:",inline"`
	AnimeTitle   string    `json:"anime_title"`
	AnimeImage   string    `json:"anime_image"`
	AnimeType    string    `json:"anime_type"`
	AnimeEpisodes int      `json:"anime_episodes"`
	AnimeStatus  string    `json:"anime_status"`
	AnimeScore   float64   `json:"anime_score"`
}

type AnimeStats struct {
	TotalWatched     int `json:"total_watched"`
	TotalPlanToWatch int `json:"total_plan_to_watch"`
	TotalWatching    int `json:"total_watching"`
	TotalWaiting     int `json:"total_waiting"`
	TotalEpisodes    int `json:"total_episodes"`
	AverageRating    float64 `json:"average_rating"`
}