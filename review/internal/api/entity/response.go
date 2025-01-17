package entity

type ApplyCouponRes struct {
	Value                 int
	AppliedDiscount       int
	ApplicationSuccessful bool
}

type ListCouponsRes struct {
	Coupons []Coupon
}

type Coupon struct {
	ID             string
	Code           string
	Discount       int
	MinBasketValue int
}
