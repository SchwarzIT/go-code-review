package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"coupon_service/internal/api/entity"
)

func (a *API) Apply(c *gin.Context) {
	code := c.Param("code")
	if len(code) == 0 {
		c.Status(http.StatusBadRequest)
		log.Error().Msg("Empty coupon code")
		return
	}
	apiReq := entity.ApplyCouponReq{}
	if err := c.ShouldBindJSON(&apiReq); err != nil {
		c.Status(http.StatusBadRequest)
		log.Error().Err(err).Msg("Failed to bind JSON")
		return
	}
	basket, err := a.svc.ApplyCoupon(apiReq.Value, apiReq.AppliedDiscount, code)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		log.Error().Err(err).Msg("Failed to apply coupon")
		return
	}
	c.JSON(http.StatusOK, basket)
}

func (a *API) Create(c *gin.Context) {
	apiReq := entity.CreateCouponReq{}
	if err := c.ShouldBindJSON(&apiReq); err != nil {
		c.Status(http.StatusBadRequest)
		log.Error().Err(err).Msg("Failed to bind JSON")
		return
	}
	err := a.svc.CreateCoupon(apiReq.Discount, apiReq.Code, apiReq.MinBasketValue)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		log.Error().Err(err).Msg("Failed to create coupon")
		return
	}
	c.Status(http.StatusCreated)
}

func (a *API) List(c *gin.Context) {
	// extract codes as filter from query string
	codesParam := c.Query("codes")
	codes := []string{}
	if codesParam != "" {
		codes = strings.Split(codesParam, ",")
	}
	coupons, err := a.svc.ListCoupons(codes...)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		log.Error().Err(err).Msg("Failed to list coupons")
		return
	}
	c.JSON(http.StatusOK, coupons)
}
