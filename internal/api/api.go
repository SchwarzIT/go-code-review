package api

import (
	"context"
	"coupon_service/internal/mytypes"
	"coupon_service/internal/service/entity"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Service interface {
	ApplyCoupon(entity.Basket, string) (*entity.Basket, error)
	CreateCoupon(int, string, int) (string, error)
	GetCoupons([]string) ([]entity.Coupon, error)
}

// Config store main api settings
// *Fields are with camel case to be read by the godotenv
type Config struct {
	PORT             string               `env:"API_PORT"`
	ENV              mytypes.Environment  `env:"API_ENV"`
	TIME_ALIVE       mytypes.MyDuration   `env:"API_TIME_ALIVE"`
	SHUTDOWN_TIMEOUT mytypes.MyDuration   `env:"API_SHUTDOWN_TIMEOUT"`
	ALLOW_ORIGINS    mytypes.AllowOrigins `env:"API_ALLOW_ORIGINS"`
}

type API struct {
	srv *http.Server
	mux *gin.Engine
	svc Service
	cfg Config
}

func New[T Service](cfg Config, svc T) (*API, error) {
	var logger *zap.Logger
	var err error

	if cfg.ENV == mytypes.Production {
		log.Println("Running in production mode")
		gin.SetMode(gin.ReleaseMode)

		logger, err = zap.NewProduction()
		if err != nil {
			return nil, err
		}
	} else {
		log.Println("Running in Development mode")
		gin.SetMode(gin.DebugMode)

		logger, err = zap.NewDevelopment()
		if err != nil {
			return nil, err
		}
	}

	r := initializeGinEngine(cfg, logger)
	api := &API{
		mux: r,
		cfg: cfg,
		svc: svc,
	}
	return api.withServer().withRoutes(), nil
}

func initializeGinEngine(cfg Config, logger *zap.Logger) *gin.Engine {
	var router *gin.Engine

	if cfg.ENV == mytypes.Production {
		router = gin.New()
		router.Use(ginLogger(logger), gin.Recovery())
		router.Use(cors.New(cors.Config{
			AllowOrigins:     cfg.ALLOW_ORIGINS,
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
			AllowCredentials: true,
			MaxAge:           time.Duration(300) * time.Second,
		}))
	} else {
		router = gin.Default()
		router.Use(cors.Default())
	}

	return router
}

func ginLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		end := time.Now()
		latency := end.Sub(start)

		if len(c.Errors) > 0 {
			for _, e := range c.Errors.Errors() {
				logger.Error(e)
			}
		} else {
			logger.Info("Incoming request",
				zap.String("method", c.Request.Method),
				zap.String("path", path),
				zap.String("query", query),
				zap.Int("status", c.Writer.Status()),
				zap.Duration("latency", latency),
				zap.String("client_ip", c.ClientIP()),
			)
		}
	}
}

func (a *API) withServer() *API {
	a.srv = &http.Server{
		Addr:              fmt.Sprintf(":%s", a.cfg.PORT),
		Handler:           a.mux,
		ReadHeaderTimeout: time.Duration(5) * time.Second,
		ReadTimeout:       time.Duration(10) * time.Second,
		WriteTimeout:      time.Duration(5) * time.Second,
		IdleTimeout:       time.Duration(10) * time.Second,
	}
	return a
}

func (a *API) withRoutes() *API {
	apiGroup := a.mux.Group("/api")
	apiGroup.POST("/apply", a.Apply)
	apiGroup.POST("/create", a.Create)
	apiGroup.GET("/coupons", a.Get)
	return a
}

func (a *API) Start() (err error) {
	err = a.srv.ListenAndServe()
	if err != nil {
		return
	}
	return
}

func (a *API) Shutdown(ctx context.Context) error {
	return a.srv.Shutdown(ctx)
}
