package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"coupon_service/internal/api"
	"coupon_service/internal/config"
	"coupon_service/internal/repository/memdb"
	"coupon_service/internal/service"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	cfg, err := config.New()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to configure app")
	}

	repo := memdb.New()
	svc := service.New(repo)
	app := api.New(cfg.Host, cfg.Port, svc)

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
