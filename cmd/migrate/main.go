package main

import (
	"database/sql"

	"github.com/FreakyGranny/anti-brute-force/internal/logger"
	"github.com/FreakyGranny/anti-brute-force/internal/storage"
	_ "github.com/FreakyGranny/anti-brute-force/migrations"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/pressly/goose"
	"github.com/rs/zerolog/log"
)

func main() {
	cfg, err := NewConfig()
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("unable initialize config")
	}
	if err := logger.SetLogLevel(cfg.Logger.Level); err != nil {
		log.Fatal().
			Err(err).
			Msg("unable to initialize logger")
	}
	dsn := storage.BuildDsn(
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.Username,
		cfg.DB.Password,
		"",
		cfg.DB.SSLEnable,
	)

	err = goose.SetDialect("postgres")
	if err != nil {
		log.Fatal().Err(err).Msg("")
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to connect")
	}
	defer db.Close()

	if err := goose.Run("up", db, "./"); err != nil {
		log.Error().Err(err).Msg("goose run")
	}
}
