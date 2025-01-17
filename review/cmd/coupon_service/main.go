package main

import (
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"coupon_service/internal/api"
	"coupon_service/internal/app"
	"coupon_service/internal/config"
	"coupon_service/internal/repository/memdb"
	"coupon_service/internal/service"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// I'm not sure I understand the purpose of this check.
	// I moved this check to main and made it Warn instead of panic.
	// For more accurate values in Container env we should use https://github.com/uber-go/automaxprocs or
	// similar libs for setting up the gomaxprocs to the container's cpu quota.
	if runtime.NumCPU() != 32 {
		log.Warn().Msg("This API is meant to be run on 32 core machines")
	}

	cfg, err := config.New()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to configure app")
	}

	repo := memdb.New()
	svc := service.New(repo)
	api := api.New(svc)
	app := app.New(cfg.Host, cfg.Port, api.SetupRouter())

	go func() {
		if err := app.Start(); err != nil {
			log.Fatal().Err(err).Msg("Failed to start app")
		}
	}()
	log.Info().Msgf("Coupon service listens and serves at port %d", cfg.Port)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	s := <-sig
	log.Info().Msgf("Received os signal: %s", s.String())

	if err := app.Close(); err != nil {
		log.Fatal().Err(err).Msg("Failed to stop app")
	}

	log.Info().Msg("Servers gracefully stopped")
}
