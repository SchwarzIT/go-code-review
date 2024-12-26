package entity

type Coupon struct {
	ID             string
	Code           string
	Discount       int
	MinBasketValue int
}
