package api

import (
	. "coupon_service/internal/api/entity"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func (a *API) Apply(c *gin.Context) {
	apiReq := ApplicationRequest{}
	if err := c.ShouldBindJSON(&apiReq); err != nil {
		return
	}
	err := a.svc.ApplyCoupon(&apiReq.Basket, apiReq.Code)
	if err != nil {
		SendError(c, "error to apply the coupon", err.Error(), http.StatusBadRequest)
		return
	}
	SendSuccess(c, "coupon applied successfully", apiReq.Basket)
}

func (a *API) Create(c *gin.Context) {
	apiReq := CouponRequest{}
	if err := c.ShouldBindJSON(&apiReq); err != nil {
		SendError(c, "error to read the body", err.Error(), http.StatusBadRequest)
		return
	}
	coupon, err := a.svc.CreateCoupon(apiReq.Discount, apiReq.Code, apiReq.MinBasketValue)
	if err != nil {
		SendError(c, "error to create coupon", err.Error(), http.StatusBadRequest)
		return
	}
	SendSuccess(c, fmt.Sprintf("coupon %s created successfully", coupon.ID), coupon)
}

func (a *API) Get(c *gin.Context) {
	codes := c.Query("codes")
	if codes == "" {
		SendError(c, "query parameter 'codes' is required", "", http.StatusBadRequest)
		return
	}
	codeList := strings.Split(codes, ",")
	if len(codeList) == 0 {
		SendError(c, "query parameter 'codes' cannot be empty", "", http.StatusBadRequest)
		return
	}

	coupons, err := a.svc.GetCoupons(codeList)
	if err != nil {
		SendError(c, "error to retrieve all coupons", err.Error(), http.StatusBadRequest)
		return
	}
	SendSuccess(c, "found all coupons", coupons)
}
