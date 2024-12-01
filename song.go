package song

type Song struct {
	Id          int    `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	Artist      string `json:"artist" db:"artist"`
	ReleaseDate string `json:"release_date" json:"release_date"`
	Text        string `json:"text" db:"text"`
	Link        string `json:"link" db:"link"`
}

type SongInput struct {
	Id          int     `json:"id" db:"id"`
	Name        *string `json:"name" db:"name"`
	Artist      *string `json:"artist" db:"artist"`
	ReleaseDate *string `json:"release_date" db:"release_date"`
	Text        *string `json:"text" db:"text"`
	Link        *string `json:"link" db:"link"`
}
