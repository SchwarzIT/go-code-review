package api

import (
	"context"

	"github.com/gin-gonic/gin"

	"coupon_service/internal/service/entity"
)

type Service interface {
	ApplyCoupon(context.Context, int, int, string) (entity.Basket, error)
	CreateCoupon(context.Context, int, string, int) error
	ListCoupons(context.Context, ...string) ([]entity.Coupon, error)
}

type API struct {
	svc Service
}

func New(svc Service) *API {
	return &API{
		svc: svc,
	}
}

func (a *API) SetupRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(gin.Recovery())

	group := r.Group("/api")
	group.POST("/coupons/:code/basket", a.Apply)
	group.POST("/coupons", a.Create)
	group.GET("/coupons", a.List)

	return r
}
