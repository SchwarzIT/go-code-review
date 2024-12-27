package main

import (
	"context"
	"coupon_service/internal/api"
	"coupon_service/internal/config"
	"coupon_service/internal/repository/memdb"
	"coupon_service/internal/service"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var interruptSignals = []os.Signal{
	os.Interrupt,
	syscall.SIGTERM,
	syscall.SIGINT,
}

func main() {
	cfg, err := config.NewDefault("")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	repo := memdb.NewRepository()
	svc := service.New(repo)

	server, err := api.New(cfg.API, svc)
	if err != nil {
		log.Fatalf("Failed to setup the server: %v", err)
	}

	serverErrors := make(chan error, 1)

	go func() {
		log.Printf("Starting Coupon service server on port %s", cfg.API.PORT)
		serverErrors <- server.Start()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, interruptSignals...)

	waitForShutdown(serverErrors, quit, cfg.API.TIME_ALIVE.ParseTimeDuration())

	ctx, cancel := context.WithTimeout(context.Background(), cfg.API.SHUTDOWN_TIMEOUT.ParseTimeDuration())
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Graceful shutdown failed: %v", err)
	}

}

func waitForShutdown(serverErrors <-chan error, quit <-chan os.Signal, TIME_ALIVE time.Duration) {
	if TIME_ALIVE > 0 {
		log.Printf("Server alive until %v", time.Now().Add(TIME_ALIVE))
		select {
		case err := <-serverErrors:
			log.Panicf("Could not start server: %v", err)
		case sig := <-quit:
			log.Printf("Received signal %s. Initiating graceful shutdown...", sig)
		case <-time.After(TIME_ALIVE):
			log.Printf("Timeout reached. Initiating graceful shutdown...")
		}
	} else {
		select {
		case err := <-serverErrors:
			log.Panicf("Could not start server: %v", err)
		case sig := <-quit:
			log.Printf("Received signal %s. Initiating graceful shutdown...", sig)
		}
	}
}
