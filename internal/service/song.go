package service

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io"
	song "music-lib"
	"music-lib/internal/repository"
	"net/http"
)

type SongService struct {
	repo         repository.Song
	musicInfoUrl string
}

func NewSongService(repo repository.Song, musicInfoUrl string) *SongService {
	return &SongService{
		repo:         repo,
		musicInfoUrl: musicInfoUrl,
	}
}

type MusicInfoResponse struct {
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

func GetMusicInfo(url, group, song string) MusicInfoResponse {
	urlInfo := url + "/info"

	req, err := http.NewRequest("GET", urlInfo, nil)
	if err != nil {
		logrus.Error(err)
		return MusicInfoResponse{}
	}

	q := req.URL.Query()
	q.Add("group", group)
	q.Add("song", song)
	req.URL.RawQuery = q.Encode()

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		logrus.Error(err)
		return MusicInfoResponse{}
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		logrus.Error(err)
		return MusicInfoResponse{}
	}

	var musicInfo MusicInfoResponse

	json.Unmarshal(body, &musicInfo)

	return musicInfo
}

func (s SongService) Add(song song.Song) (int, error) {
	musicInfo := GetMusicInfo(s.musicInfoUrl, song.Artist, song.Name)
	song.ReleaseDate = musicInfo.ReleaseDate
	song.Text = musicInfo.Text
	song.Link = musicInfo.Link

	return s.repo.Add(song)
}

func (s SongService) GetById(id int) (song.SongInput, error) {
	return s.repo.GetById(id)
}

func (s SongService) GetAll(songData song.Song, limit, page int) ([]song.SongInput, error) {
	return s.repo.GetAll(songData, limit, page)
}

func (s SongService) Update(id int, song song.SongInput) error {
	return s.repo.Update(id, song)
}

func (s SongService) Delete(id int) error {
	return s.repo.Delete(id)
}

func (s SongService) GetSongTextById(id int) (string, error) {
	return s.repo.GetSongTextById(id)
}
