package main

import (
	"context"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	song "music-lib"
	"music-lib/internal/handler"
	"music-lib/internal/repository"
	"music-lib/internal/service"
	"os"
	"os/signal"
)

// @title Music Lib API
// @version 1.0
// @host localhost:8000
// @Basepath /

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	// config
	err := godotenv.Load()
	if err != nil {
		logrus.Fatal("Error loading .env file")
	}

	// repository
	db, err := repository.NewPostgresDB(repository.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSL_MODE"),
	})
	if err != nil {
		logrus.Fatalf("Error connecting to database: %v", err.Error())
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos, os.Getenv("MUSIC_INFO_URL"))
	handlers := handler.NewHandler(services)

	srv := new(song.Server)
	go func() {
		if err := srv.Run(os.Getenv("APP_PORT"), handlers.Init()); err != nil {
			logrus.Fatalf("Error starting server: %v", err.Error())
		}
	}()

	logrus.Print("MusicLib started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	logrus.Print("MusicLib stopping")

	if err := srv.Stop(context.Background()); err != nil {
		logrus.Fatalf("Error stopping server: %v", err.Error())
	}

	if err := db.Close(); err != nil {
		logrus.Fatalf("Error closing database: %v", err.Error())
	}
}
