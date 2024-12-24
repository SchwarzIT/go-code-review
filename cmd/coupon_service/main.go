package main

import (
	"coupon_service/internal/api"
	"coupon_service/internal/config"
	"coupon_service/internal/repository/memdb"
	"coupon_service/internal/service"
	"fmt"
	"time"
)

var (
	cfg  = config.New()
	repo = memdb.New()
)

func main() {
	svc := service.New(repo)
	fmt.Println(cfg)
	server := api.New(cfg.API, svc)
	server.Start()
	fmt.Println("Starting Coupon service server")
	<-time.After(1 * time.Hour * 24 * 365)
	fmt.Println("Coupon service server alive for a year, closing")
	server.Close()
}
