package entity

type CouponRequest struct {
	Discount       int    `json:"discount" binding:"required"`
	Code           string `json:"code" binding:"required"`
	MinBasketValue int    `json:"min_basket_value" binding:"required"`
}
