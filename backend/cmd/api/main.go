package main

import (
	"os"

	"github.com/uniflou/uniflou-lms/backend/internal/config"
	"github.com/uniflou/uniflou-lms/backend/internal/server"
	"github.com/uniflou/uniflou-lms/backend/pkg/database"
	"github.com/uniflou/uniflou-lms/backend/pkg/logger"
)

func main() {
	log := logger.New()
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("config")
	}
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("database")
	}
	srv := server.New(cfg, db, log)
	if err := srv.Run(); err != nil {
		log.Error().Err(err).Msg("server stopped")
		os.Exit(1)
	}
}
