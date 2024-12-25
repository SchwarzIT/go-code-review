package api

import (
	. "coupon_service/internal/api/entity"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func (a *API) Apply(c *gin.Context) {
	apiReq := ApplicationRequest{}
	if err := c.ShouldBindJSON(&apiReq); err != nil {
		return
	}
	basket, err := a.svc.ApplyCoupon(apiReq.Basket, apiReq.Code)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, basket)
}

func (a *API) Create(c *gin.Context) {
	apiReq := Coupon{}
	if err := c.ShouldBindJSON(&apiReq); err != nil {
		return
	}
	id, err := a.svc.CreateCoupon(apiReq.Discount, apiReq.Code, apiReq.MinBasketValue)
	if err != nil {
		return
	}
	c.Status(http.StatusOK)
	c.Writer.Write([]byte(id))
}

func (a *API) Get(c *gin.Context) {
	codes := c.Query("codes")
	if codes == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter 'codes' is required"})
		return
	}
	codeList := strings.Split(codes, ",")
	if len(codeList) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter 'codes' cannot be empty"})
		return
	}

	coupons, err := a.svc.GetCoupons(codeList)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "failed not found coupons"})
		return
	}
	c.JSON(http.StatusOK, coupons)
}
