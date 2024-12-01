package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	song "music-lib"
	"strconv"
	"strings"
)

type SongPostgres struct {
	db *sql.DB
}

func NewSongPostgres(db *sql.DB) *SongPostgres {
	return &SongPostgres{db: db}
}

func (r *SongPostgres) Add(song song.Song) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		logrus.Error(err)
		return 0, err
	}
	defer tx.Rollback()

	var id int
	cols := make([]string, 0)
	values := make([]string, 0)
	args := make([]interface{}, 0)

	if song.Name != "" {
		cols = append(cols, "name")
		args = append(args, song.Name)
		values = append(values, "$"+strconv.Itoa(len(cols)))
	}
	if song.Artist != "" {
		cols = append(cols, "artist")
		args = append(args, song.Artist)
		values = append(values, "$"+strconv.Itoa(len(cols)))
	}
	if song.ReleaseDate != "" {
		cols = append(cols, "release_date")
		args = append(args, song.ReleaseDate)
		values = append(values, "TO_DATE($"+strconv.Itoa(len(cols))+", 'DD.MM.YYYY')")
	}
	if song.Link != "" {
		cols = append(cols, "link")
		args = append(args, song.Link)
		values = append(values, "$"+strconv.Itoa(len(cols)))
	}
	if song.Text != "" {
		cols = append(cols, "text")
		args = append(args, song.Text)
		values = append(values, "$"+strconv.Itoa(len(cols)))
	}

	query := fmt.Sprintf("INSERT INTO songs (%v) VALUES (%v) RETURNING id", strings.Join(cols, ", "), strings.Join(values, ", "))
	row := tx.QueryRow(query, args...)

	if err := row.Scan(&id); err != nil {
		logrus.Error(err)
		return 0, err
	}

	return id, tx.Commit()
}

func (r *SongPostgres) GetById(id int) (song.SongInput, error) {
	var songRes song.SongInput
	query := "SELECT id, name, artist, TO_CHAR(release_date, 'DD.MM.YYYY'), text, link FROM songs WHERE id = $1"
	if err := r.db.QueryRow(query, id).Scan(&songRes.Id, &songRes.Name, &songRes.Artist, &songRes.ReleaseDate, &songRes.Text, &songRes.Link); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return song.SongInput{}, nil
		}
		logrus.Error(err)
		return song.SongInput{}, err
	}
	return songRes, nil
}

func (r *SongPostgres) GetAll(songData song.Song, limit, page int) ([]song.SongInput, error) {
	whereArr := make([]string, 0)
	args := make([]interface{}, 0)

	if songData.Name != "" {
		whereArr = append(whereArr, "LOWER(name) LIKE LOWER('%' || $"+strconv.Itoa(len(whereArr)+1)+" || '%')")
		args = append(args, songData.Name)
	}
	if songData.Artist != "" {
		whereArr = append(whereArr, "LOWER(artist) LIKE LOWER('%' || $"+strconv.Itoa(len(whereArr)+1)+" || '%')")
		args = append(args, songData.Artist)
	}
	if songData.ReleaseDate != "" {
		whereArr = append(whereArr, "release_date = TO_DATE($"+strconv.Itoa(len(whereArr)+1)+", 'DD.MM.YYYY')")
		args = append(args, songData.ReleaseDate)
	}
	if songData.Text != "" {
		whereArr = append(whereArr, "LOWER(text) LIKE LOWER('%' || $"+strconv.Itoa(len(whereArr)+1)+" || '%')")
		args = append(args, songData.Text)
	}
	if songData.Link != "" {
		whereArr = append(whereArr, "LOWER(link) LIKE LOWER('%' || $"+strconv.Itoa(len(whereArr)+1)+" || '%')")
		args = append(args, songData.Link)
	}

	whereStr := ""
	if len(whereArr) > 0 {
		whereStr = "WHERE " + strings.Join(whereArr, " AND ")
	}

	offset := limit*page - limit

	query := fmt.Sprintf("SELECT id, name, artist, TO_CHAR(release_date, 'DD.MM.YYYY'), text, link FROM songs %v LIMIT %v OFFSET %v", whereStr, limit, offset)
	rows, err := r.db.Query(query, args...)
	if err != nil {
		logrus.Error(err)
		return []song.SongInput{}, err
	}
	defer rows.Close()

	var songs []song.SongInput

	for rows.Next() {
		var songRes song.SongInput
		if err := rows.Scan(&songRes.Id, &songRes.Name, &songRes.Artist, &songRes.ReleaseDate, &songRes.Text, &songRes.Link); err != nil {
			logrus.Error(err)
			return []song.SongInput{}, err
		}
		songs = append(songs, songRes)
	}
	if err := rows.Err(); err != nil {
		logrus.Error(err)
		return []song.SongInput{}, err
	}
	return songs, nil
}

func (r *SongPostgres) Update(id int, song song.SongInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)

	if song.Name != nil {
		setValues = append(setValues, "name = $"+strconv.Itoa(len(args)+1))
		args = append(args, *song.Name)
	}
	if song.Artist != nil {
		setValues = append(setValues, "artist = $"+strconv.Itoa(len(args)+1))
		args = append(args, *song.Artist)
	}
	if song.ReleaseDate != nil {
		setValues = append(setValues, "release_date = TO_DATE($"+strconv.Itoa(len(args)+1)+", 'DD.MM.YYYY')")
		args = append(args, *song.ReleaseDate)
	}
	if song.Link != nil {
		setValues = append(setValues, "link = $"+strconv.Itoa(len(args)+1))
		args = append(args, *song.Link)
	}
	if song.Text != nil {
		setValues = append(setValues, "text = $"+strconv.Itoa(len(args)+1))
		args = append(args, *song.Text)
	}

	if len(args) == 0 {
		return errors.New("no arguments passed to Update")
	}

	tx, err := r.db.Begin()
	if err != nil {
		logrus.Error(err)
		return err
	}
	defer tx.Rollback()

	query := fmt.Sprintf("UPDATE songs SET %v WHERE id = $%v", strings.Join(setValues, ", "), len(args)+1)

	_, err = tx.Exec(query, append(args, id)...)
	if err != nil {
		logrus.Error(err)
		return err
	}

	if err = tx.Commit(); err != nil {
		logrus.Error(err)
		return err
	}

	return nil
}

func (r *SongPostgres) Delete(id int) error {
	query := "DELETE FROM songs WHERE id = $1"
	_, err := r.db.Exec(query, id)
	if err != nil {
		logrus.Error(err)
	}
	return err
}

func (r *SongPostgres) GetSongTextById(id int) (string, error) {
	var text string
	query := "SELECT text FROM songs WHERE id = $1"
	if err := r.db.QueryRow(query, id).Scan(&text); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}
		logrus.Error(err)
		return "", err
	}
	return text, nil
}
