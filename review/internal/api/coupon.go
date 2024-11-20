package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"coupon_service/internal/entity"
)

func (a *API) Apply(c *gin.Context) {
	apiReq := entity.ApplicationRequest{}
	if err := c.ShouldBindJSON(&apiReq); err != nil {
		return
	}

	err := a.svc.ApplyCoupon(&apiReq.Basket, apiReq.Code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid coupon",
			"message": fmt.Sprintf("The coupon code %s is invalid", apiReq.Code),
		})
		return
	}

	c.JSON(http.StatusOK, apiReq.Basket)
}

func (a *API) Create(c *gin.Context) {
	apiReq := entity.Coupon{}
	if err := c.ShouldBindJSON(&apiReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid format",
			"message": "The request body contains invalid format",
		})
		return
	}

	err := a.svc.CreateCoupon(apiReq.Discount, apiReq.Code, apiReq.MinBasketValue)
	if err != nil {
		//TODO: create custom error for different type of error responses
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   fmt.Sprintf("%v", err),
			"message": "An error was found while creating the coupon",
		})
		return
	}

	coupon, err := a.svc.GetCoupon(apiReq.Code)
	c.JSON(http.StatusOK, gin.H{"data": coupon})
}

func (a *API) Get(c *gin.Context) {
	apiReq := entity.CouponRequest{
		Codes: []string{},
	}
	if err := c.ShouldBindJSON(&apiReq); err != nil {
		log.Printf("error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid format",
			"message": "The requestbody contains invalid format",
		})
		return
	}

	coupons := a.svc.GetCoupons(apiReq.Codes)
	response := map[string]interface{}{
		"coupons": coupons,
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}
