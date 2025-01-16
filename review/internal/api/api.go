package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"coupon_service/internal/service/entity"
)

type Service interface {
	ApplyCoupon(int, int, string) (entity.Basket, error)
	CreateCoupon(int, string, int) error
	ListCoupons(...string) ([]entity.Coupon, error)
}

type API struct {
	srv *http.Server
	svc Service
}

func New(host string, port int, svc Service) *API {
	api := &API{
		svc: svc,
	}

	// TODO: rework
	r := api.setupRouter()

	api.srv = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", host, port),
		Handler: r,
	}

	return api
}

func (a *API) setupRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(gin.Recovery())

	group := r.Group("/api")
	group.POST("/coupons/:code/basket", a.Apply)
	group.POST("/coupons", a.Create)
	group.GET("/coupons", a.List)

	return r
}

func (a *API) Start() error {
	if err := a.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to serve http: %w", err)
	}
	return nil
}

func (a *API) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown http server: %w", err)
	}
	return nil
}
