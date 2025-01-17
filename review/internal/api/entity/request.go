package entity

type ApplyCouponReq struct {
	Value           int
	AppliedDiscount int
}

type CreateCouponReq struct {
	Discount       int
	Code           string
	MinBasketValue int
}
