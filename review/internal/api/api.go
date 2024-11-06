package api

import (
	"context"
	"coupon_service/internal/service/entity"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type couponService interface {
	ApplyCoupon(basket entity.Basket, code string) (*entity.Basket, error)
	CreateCoupon(discount int, code string, minBasketValue int) error
	GetCoupons(codes []string) ([]entity.Coupon, error)
}

// Config is the configuration for the API
type Config struct {
	Host string
	Port int
}

// API is the API for the service
type API struct {
	srv    *http.Server
	mux    *gin.Engine
	svc    couponService
	config Config
}

// New creates a new API
func New(cfg Config, svc couponService) API {
	gin.SetMode(gin.ReleaseMode)
	r := new(gin.Engine)
	r = gin.New()
	r.Use(gin.Recovery())

	api := API{
		mux:    r,
		svc:    svc,
		config: cfg,
	}
	return api.withServer().withRoutes()
}

func (a API) withServer() API {
	ch := make(chan API)
	go func() {
		a.srv = &http.Server{
			Addr:    fmt.Sprintf(":%d", a.config.Port),
			Handler: a.mux,
		}
		ch <- a
	}()

	return <-ch
}

func (a API) withRoutes() API {
	apiGroup := a.mux.Group("/api")
	apiGroup.POST("/apply", a.Apply)
	apiGroup.POST("/coupons", a.Create)
	apiGroup.GET("/coupons", a.Get)
	return a
}

// Start starts the API
func (a API) Start() {
	if err := a.srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

// Close closes the API
func (a API) Close() {
	<-time.After(5 * time.Second)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.srv.Shutdown(ctx); err != nil {
		log.Println(err)
	}
}
