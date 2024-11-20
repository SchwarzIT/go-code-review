package entity

// TODO: write validations for Discount and MinBasketValue

type Coupon struct {
	ID             string `json:"id,omitempty"`
	Code           string `json:"code,omitempty"`
	Discount       int    `json:"discount,omitempty"`
	MinBasketValue int    `json:"min_basket_value,omitempty"`
}
