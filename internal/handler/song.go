package handler

import (
	"github.com/gin-gonic/gin"
	song "music-lib"
	"net/http"
	"strconv"
)

const (
	LimitDefault = 20
	PageDefault  = 1
)

// @Summary Add Song
// @Tags songs
// @Description add song to database
// @ID add-song
// @Accept json
// @Produce json
// @Param song body song.Song true "song info"
// @Success 200 {integer} integer "New song ID"
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/songs [post]
func (h *Handler) addSong(c *gin.Context) {
	var input song.Song
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.Song.Add(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id})
}

// @Summary Get All Songs
// @Tags songs
// @Description get all songs
// @ID get-all-songs
// @Accept json
// @Produce json
// @Param name query string false "song name"
// @Param artist query string false "song artist"
// @Param release_date query string false "song release date"
// @Param text query string false "song text"
// @Param link query string false "song link"
// @Param limit query int false "limit songs" default(LimitDefault)
// @Param page query int false "selected page" default(PageDefault)
// @Success 200 {object} []song.SongInput
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/songs [get]
func (h *Handler) getAllSongs(c *gin.Context) {
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil || limit <= 0 {
		limit = LimitDefault
	}

	page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page <= 0 {
		page = PageDefault
	}

	var songData song.Song

	songData.Name = c.Query("name")
	songData.Artist = c.Query("artist")
	songData.ReleaseDate = c.Query("release_date")
	songData.Text = c.Query("text")
	songData.Link = c.Query("link")

	items, err := h.services.Song.GetAll(songData, limit, page)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, items)
}

// @Summary Get Song By ID
// @Tags songs
// @Description get song by id
// @ID get-song-by-id
// @Accept json
// @Produce json
// @Param id path int true "Song ID"
// @Success 200 {object} song.Song
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/songs/{id} [get]
func (h *Handler) getSongById(c *gin.Context) {
	songId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	item, err := h.services.Song.GetById(songId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, item)
}

// @Summary Update Song
// @Tags songs
// @Description update song
// @ID update-song
// @Accept json
// @Produce json
// @Param id path int true "Song ID"
// @Param song body song.SongInput true "New song info"
// @Success 200 {object} statusResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/songs/{id} [put]
func (h *Handler) updateSong(c *gin.Context) {
	songId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var input song.SongInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.services.Update(songId, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{"ok"})
}

// @Summary Delete Song
// @Tags songs
// @Description delete song
// @ID delete-song
// @Accept json
// @Produce json
// @Param id path int true "Song ID"
// @Success 200 {object} song.Song
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/songs/{id} [delete]
func (h *Handler) deleteSong(c *gin.Context) {
	songId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err = h.services.Song.Delete(songId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{"ok"})
}

// @Summary Get Song Text
// @Tags songs
// @Description get song text
// @ID get-song-text
// @Accept json
// @Produce json
// @Param id path int true "Song ID"
// @Success 200 {string} string "Song text"
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/songs/{id}/text [get]
func (h *Handler) getSongText(c *gin.Context) {
	songId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	text, err := h.services.Song.GetSongTextById(songId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, text)
}
