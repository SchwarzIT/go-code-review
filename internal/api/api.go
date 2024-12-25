package api

import (
	"context"
	"coupon_service/internal/myduration"
	"coupon_service/internal/service/entity"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Service interface {
	ApplyCoupon(entity.Basket, string) (*entity.Basket, error)
	CreateCoupon(int, string, int) (string, error)
	GetCoupons([]string) ([]entity.Coupon, error)
}

type Config struct {
	Port      string                `env:"API_PORT"`
	Env       string                `env:"API_ENV"`
	TimeAlive myduration.MyDuration `env:"API_TIMEALIVE"`
}

type API struct {
	srv *http.Server
	MUX *gin.Engine
	svc Service
	CFG Config
}

func New[T Service](cfg Config, svc T) API {
	r := gin.Default()
	return API{
		MUX: r,
		CFG: cfg,
		svc: svc,
	}.withServer().withRoutes()
}

func (a API) withServer() API {
	a.srv = &http.Server{
		Addr:    fmt.Sprintf(":%s", a.CFG.Port),
		Handler: a.MUX,
	}
	return a
}

func (a API) withRoutes() API {
	apiGroup := a.MUX.Group("/api")
	apiGroup.POST("/apply", a.Apply)
	apiGroup.POST("/create", a.Create)
	apiGroup.GET("/coupons", a.Get)
	return a
}

func (a API) Start() (err error) {
	err = a.srv.ListenAndServe()
	if err != nil {
		return
	}
	return
}

func (a API) Close() {
	<-time.After(5 * time.Second)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.srv.Shutdown(ctx); err != nil {
		log.Println(err)
	}
}
