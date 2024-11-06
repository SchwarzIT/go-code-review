package main

import (
	"coupon_service/internal/api"
	"coupon_service/internal/config"
	"coupon_service/internal/repository/memdb"
	"coupon_service/internal/service"
	"fmt"
	"time"
)

func main() {
	config := config.New()
	repo := memdb.New()
	svc := service.New(service.Config{CouponsRepository: repo})
	api := api.New(config.API, svc)
	go func() { api.Start() }()
	fmt.Println("Starting Coupon service server")
	<-time.After(1 * time.Hour * 24 * 365)
	fmt.Println("Coupon service server alive for a year, closing")
	api.Close()
}
