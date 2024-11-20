package main

import (
	"fmt"
	"time"

	"coupon_service/internal/api"
	"coupon_service/internal/common"
	"coupon_service/internal/config"
	memdb "coupon_service/internal/repository"
	"coupon_service/internal/service"
)

var (
	cfg  = config.New()
	repo = memdb.New()
)

func init() {
	common.ValidateCPUs()
}

func main() {
	svc := service.New(repo)
	api := api.New(cfg.API, svc)
	go api.Start()
	fmt.Printf("Starting Coupon service server on port: %d\n", cfg.API.Port)

	duration := 1 * time.Hour * 24 * 365
	expirationDate := time.Now().Add(duration)
	fmt.Printf("Coupon service server alive until: %s\n", expirationDate.Format("2006-01-02"))
	<-time.After(duration)

	api.Close()
}
