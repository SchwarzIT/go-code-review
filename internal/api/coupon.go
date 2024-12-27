package api

import (
	"coupon_service/internal/service"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type ApplicationRequest struct {
	Code   string         `json:"code" binding:"required"`
	Basket service.Basket `json:"basket" binding:"required"`
}

type CouponRequest struct {
	Discount       int    `json:"discount" binding:"required"`
	Code           string `json:"code" binding:"required"`
	MinBasketValue int    `json:"min_basket_value" binding:"required"`
}

// Apply handles the HTTP POST request to apply a coupon to a basket.
//
// @Summary      Apply coupon
// @Description  Applies a coupon code to the given basket
// @Tags         Coupons
// @Accept       json
// @Produce      json
// @Param        body  body      ApplicationRequest  true  "Basket and coupon code"
// @Success      200   {object}  service.Basket
// @Failure      400   {string}  string
// @Failure      500   {string}  string
// @Router       /apply [post]
func (a *API) Apply(c *gin.Context) {
	apiReq := ApplicationRequest{}
	if err := c.ShouldBindJSON(&apiReq); err != nil {
		SendError(c, "error to read the body", err.Error(), http.StatusBadRequest)
		return
	}
	err := a.svc.ApplyCoupon(&apiReq.Basket, apiReq.Code)
	if err != nil {
		SendError(c, "error to apply the coupon", err.Error(), http.StatusBadRequest)
		return
	}
	SendSuccess(c, "coupon applied successfully", apiReq.Basket)
}

// Create handles the HTTP POST request for creating a new coupon.
//
// @Summary      Create a new coupon
// @Description  Creates a new coupon using the given details in the request body
// @Tags         Coupons
// @Accept       json
// @Produce      json
// @Param        body  body      CouponRequest  true  "Coupon creation data"
// @Success      201   {object}  memdb.Coupon
// @Failure      400   {string}  string
// @Failure      500   {string}  string
// @Router       /create [post]
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

// Get handles the HTTP GET request for retrieving coupons.
//
// @Summary      Get coupons by code
// @Description  Fetches multiple coupons based on a comma-separated list of codes
// @Tags         Coupons
// @Accept       json
// @Produce      json
// @Param        codes  query     string  true  "Comma-separated coupon codes (e.g. 'CODE123,CODE456')"
// @Success      200    {array}   memdb.Coupon
// @Failure      400    {string}  string
// @Failure      404    {string}  string
// @Router       /coupons [get]
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
		SendError(c, "error to found all coupons", err.Error(), http.StatusNotFound)
		return
	}
	SendSuccess(c, "found all coupons", coupons)
}
