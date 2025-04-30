package models

type Anime struct {
	MALId         int64   `json:"mal_id"`
	Title         string  `json:"title"`
	TitleEnglish  string  `json:"title_english"`
	TitleJapanese string  `json:"title_japanese"`
	Synopsis      string  `json:"synopsis"`
	ImageURL      string  `json:"image_url"`
	Type          string  `json:"type"`
	Source        string  `json:"source"`
	Episodes      int     `json:"episodes"`
	Status        string  `json:"status"`
	Airing        bool    `json:"airing"`
	Score         float64 `json:"score"`
	Rank          int     `json:"rank"`
	Popularity    int     `json:"popularity"`
	Genres        []Genre `json:"genres"`
}

type Genre struct {
	ID  int64  `json:"id"`
	Name string `json:"name"`
}

type AnimeSeatch struct {
	Query     string `json:"query" form:"query"`
	Type      string `json:"type" form:"type"`
	Status    string `json:"status" form:"status"`
	Genre     string `json:"genre" form:"genre"`
	Page      int    `json:"page" form:"page"`
	Limit     int    `json:"limit" form:"limit"`
	Sort      string `json:"sort" form:"sort"`
	SortOrder string `json:"sort_order" form:"sort_order"`
}

type AnimeList struct {
	Items      []*Anime `json:"items"`
	TotalCount int      `json:"total_count"`
	Page       int      `json:"page"`
	Limit      int      `json:"limit"`
}