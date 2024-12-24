package main

import (
	"coupon_service/internal/api"
	"coupon_service/internal/config"
	"coupon_service/internal/repository/memdb"
	"coupon_service/internal/service"
	"fmt"
	"log"
	"time"
)

var (
	repo = memdb.New()
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	svc := service.New(repo)

	server := api.New(cfg.API, svc)
	server.Start()
	fmt.Println("Starting Coupon service server")
	<-time.After(1 * time.Hour * 24 * 365)
	fmt.Println("Coupon service server alive for a year, closing")
	server.Close()
}
