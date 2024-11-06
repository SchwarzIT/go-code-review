package api

import (
	. "coupon_service/internal/api/entity"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Apply applies a coupon to a basket
func (a *API) Apply(c *gin.Context) {
	apiReq := ApplicationRequest{}
	if err := c.ShouldBindJSON(&apiReq); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	basket, err := a.svc.ApplyCoupon(apiReq.Basket, apiReq.Code)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusOK, basket)
}

// Create creates a coupon
func (a *API) Create(c *gin.Context) {
	apiReq := Coupon{}
	if err := c.ShouldBindJSON(&apiReq); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	err := a.svc.CreateCoupon(apiReq.Discount, apiReq.Code, apiReq.MinBasketValue)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	c.Status(http.StatusOK)
}

// Get gets all the coupons
func (a *API) Get(c *gin.Context) {
	apiReq := CouponRequest{}
	if err := c.ShouldBindJSON(&apiReq); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	coupons, err := a.svc.GetCoupons(apiReq.Codes)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusOK, coupons)
}
