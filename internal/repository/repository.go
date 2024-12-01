package repository

import (
	"database/sql"
	song "music-lib"
)

type Song interface {
	Add(song.Song) (int, error)
	GetById(id int) (song.SongInput, error)
	GetAll(songData song.Song, limit, page int) ([]song.SongInput, error)
	Update(id int, song song.SongInput) error
	Delete(id int) error
	GetSongTextById(id int) (string, error)
}

type Repository struct {
	Song
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Song: NewSongPostgres(db),
	}
}
