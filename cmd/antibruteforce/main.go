package main

import (
	"os"
	"os/signal"

	"github.com/FreakyGranny/anti-brute-force/internal/app"
	"github.com/FreakyGranny/anti-brute-force/internal/cache"
	"github.com/FreakyGranny/anti-brute-force/internal/logger"
	"github.com/FreakyGranny/anti-brute-force/internal/storage"

	internalhttp "github.com/FreakyGranny/anti-brute-force/internal/server/http"
	"github.com/jonboulle/clockwork"
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
	db := storage.New(storage.BuildDsn(
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.Username,
		cfg.DB.Password,
		cfg.DB.Name,
		cfg.DB.SSLEnable,
	))
	defer db.Close()

	cache, err := cache.New(cfg.Redis.Host, cfg.Redis.Port, cfg.Redis.Password)
	if err != nil {
		log.Err(err).
			Msg("redis unavailable")
		os.Exit(1)
	}
	defer cache.Close()

	a := app.New(
		db,
		cache,
		clockwork.NewRealClock(),
		cfg.Limits.User,
		cfg.Limits.Password,
		cfg.Limits.IP,
	)

	server := internalhttp.NewServer(":8090", a)

	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, os.Interrupt)

		<-signals
		signal.Stop(signals)

		if err := server.Stop(); err != nil {
			log.Err(err).
				Msg("failed to stop http server")
		}
	}()

	if err := server.Start(); err != nil {
		log.Err(err).
			Msg("failed to start http server")
		os.Exit(1)
	}

}
