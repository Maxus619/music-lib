package service

import (
	song "music-lib"
	"music-lib/internal/repository"
)

type Song interface {
	Add(song.Song) (int, error)
	GetById(id int) (song.SongInput, error)
	GetAll(songData song.Song, limit, page int) ([]song.SongInput, error)
	Update(id int, song song.SongInput) error
	Delete(id int) error
	GetSongTextById(id int) (string, error)
}

type Service struct {
	Song
}

func NewService(repos *repository.Repository, musicInfoUrl string) *Service {
	return &Service{
		Song: NewSongService(repos.Song, musicInfoUrl),
	}
}
