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
		responseError(c, http.StatusBadRequest, "Empty coupon code", nil)
		return
	}
	apiReq := entity.ApplyCouponReq{}
	if err := c.ShouldBindJSON(&apiReq); err != nil {
		responseError(c, http.StatusBadRequest, "Failed to bind JSON", err)
		return
	}
	basket, err := a.svc.ApplyCoupon(c.Request.Context(), apiReq.Value, apiReq.AppliedDiscount, code)
	if err != nil {
		responseError(c, http.StatusInternalServerError, "Failed to apply coupon", err)
		return
	}
	c.JSON(http.StatusOK, basket)
}

func (a *API) Create(c *gin.Context) {
	apiReq := entity.CreateCouponReq{}
	if err := c.ShouldBindJSON(&apiReq); err != nil {
		responseError(c, http.StatusBadRequest, "Failed to bind JSON", err)
		return
	}
	err := a.svc.CreateCoupon(c.Request.Context(), apiReq.Discount, apiReq.Code, apiReq.MinBasketValue)
	if err != nil {
		responseError(c, http.StatusInternalServerError, "Failed to create coupon", err)
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
	coupons, err := a.svc.ListCoupons(c.Request.Context(), codes...)
	if err != nil {
		responseError(c, http.StatusInternalServerError, "Failed to list coupons", err)
		return
	}
	c.JSON(http.StatusOK, coupons)
}

func responseError(c *gin.Context, code int, msg string, err error) {
	c.JSON(code, entity.Error{
		Code: code,
		Msg:  msg,
	})
	e := log.Error()
	if err != nil {
		e = e.Err(err)
	}
	e.Msg(msg)
}
