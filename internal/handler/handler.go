package handler

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "music-lib/docs"
	"music-lib/internal/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) Init() *gin.Engine {
	router := gin.New()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group("/api")
	{
		songs := api.Group("/songs")
		{
			songs.GET("/", h.getAllSongs)
			songs.GET("/:id", h.getSongById)
			songs.POST("/", h.addSong)
			songs.PUT("/:id", h.updateSong)
			songs.DELETE("/:id", h.deleteSong)

			text := songs.Group(":id/text")
			{
				text.GET("/", h.getSongText)
			}
		}
	}

	return router
}
