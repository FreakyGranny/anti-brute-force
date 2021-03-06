package main

import (
	"context"
	"net"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/FreakyGranny/anti-brute-force/internal/app"
	"github.com/FreakyGranny/anti-brute-force/internal/cache"
	"github.com/FreakyGranny/anti-brute-force/internal/logger"
	"github.com/FreakyGranny/anti-brute-force/internal/server"
	"github.com/FreakyGranny/anti-brute-force/internal/storage"
	"github.com/jonboulle/clockwork"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

func main() {
	cfg, err := NewConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("unable initialize config")
	}
	if err := logger.SetLogLevel(cfg.Logger.Level); err != nil {
		log.Fatal().Err(err).Msg("unable to initialize logger")
	}
	refreshInterval, err := time.ParseDuration(cfg.IPListRefresh)
	if err != nil {
		log.Fatal().Err(err).Msgf("wrong ip list refresh interval %s", cfg.IPListRefresh)
	}
	db, err := storage.New(storage.BuildDsn(
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.Username,
		cfg.DB.Password,
		cfg.DB.Name,
		cfg.DB.SSLEnable,
	))
	if err != nil {
		log.Fatal().Err(err).Msg("error while initializing db")
	}
	defer func() {
		if err = db.Close(); err != nil {
			log.Error().Err(err).Msg("error at close db connection")
		}
	}()

	keeper := app.NewMemIPKeeper(db)
	keeperCtx, keeperCancel := context.WithCancel(context.Background())
	if err := keeper.Refresh(keeperCtx); err != nil {
		log.Error().Err(err).Msg("can't initialize black/white lists")
		return
	}
	redisCache, err := cache.New(cfg.Redis.Host, cfg.Redis.Port, cfg.Redis.Password)
	if err != nil {
		log.Error().Err(err).Msg("redis unavailable")
		return
	}
	defer func() {
		if err = redisCache.Close(); err != nil {
			log.Error().Err(err).Msg("error at close db connection")
		}
	}()

	limiter := app.NewCacheLimiter(redisCache, clockwork.NewRealClock(), cfg.Limits.User, cfg.Limits.Password, cfg.Limits.IP)
	lsn, err := net.Listen("tcp", net.JoinHostPort(cfg.GRPC.Host, strconv.Itoa(cfg.GRPC.Port)))
	if err != nil {
		log.Error().Err(err).Msg("")
		return
	}
	srv := grpc.NewServer(grpc.UnaryInterceptor(server.LoggingInterceptor))
	service := server.New(app.New(db, keeper, limiter))
	server.RegisterABruteforceServer(srv, service)

	go keeper.Watch(keeperCtx, refreshInterval)
	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, os.Interrupt)
		<-signals
		signal.Stop(signals)
		keeperCancel()
		srv.GracefulStop()
	}()

	log.Info().Msgf("Starting server on %s", lsn.Addr().String())
	if err := srv.Serve(lsn); err != nil {
		log.Error().Err(err).Msg("failed to start grpc server")
		return
	}
}
