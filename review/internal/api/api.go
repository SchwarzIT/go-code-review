package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"coupon_service/internal/entity"
)

type Service interface {
	ApplyCoupon(*entity.Basket, string) error
	CreateCoupon(int, string, int) error
	GetCoupons([]string) []entity.Coupon
	GetCoupon(string) (*entity.Coupon, error)
}

type Config struct {
	Host string
	Port int
}

type API struct {
	srv *http.Server
	MUX *gin.Engine
	svc Service
	CFG Config
}

// New creates a new instance of the API, initializes the router, sets up routes, and configures the server.
// It accepts a configuration struct and a service that implements the Service interface.
func New[T Service](cfg Config, svc T) API {
	gin.SetMode(gin.ReleaseMode)
	r := new(gin.Engine)
	r = gin.New()
	r.Use(gin.Recovery())

	return API{
		MUX: r,
		CFG: cfg,
		svc: svc,
	}.withRoutes().withServer()
}

// withServer configures the HTTP server for the API, binding it to the host and port specified in the configuration.
// It starts the server in a separate goroutine and returns the API instance.
func (a API) withServer() API {
	ch := make(chan API)
	go func() {
		a.srv = &http.Server{
			Addr:    fmt.Sprintf("%s:%d", a.CFG.Host, a.CFG.Port),
			Handler: a.MUX,
		}
		ch <- a
	}()

	return <-ch
}

// withRoutes defines the API routes for handling requests and returns the API instance.
// It includes endpoints for applying, creating, and retrieving coupons.
func (a API) withRoutes() API {
	apiGroup := a.MUX.Group("/api")
	apiGroup.POST("/apply", a.Apply)
	apiGroup.POST("/create", a.Create)
	apiGroup.GET("/coupons", a.Get)
	return a
}

// Start begins serving HTTP requests on the specified address.
// It listens on the server's configured host and port, logging any critical errors if the server fails to start.
func (a API) Start() {
	log.Println(a.srv.Addr)
	if err := a.srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

// Close gracefully shuts down the HTTP server within a 5-second timeout period.
// This allows ongoing requests to complete before the server stops.
func (a API) Close() {
	<-time.After(5 * time.Second)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.srv.Shutdown(ctx); err != nil {
		log.Println(err)
	}
}
