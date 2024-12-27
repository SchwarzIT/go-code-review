package service

import (
	"errors"
	"fmt"
)

// Custom errors for better error handling.
var (
	ErrApplyDiscount          = errors.New("cannot apply discount to a basket of non-positive value")
	ErrCouponDiscountValue    = errors.New("discount must be a positive integer")
	ErrCouponMinBasketValue   = errors.New("minBasketValue cannot be negative")
	ErrCouponDiscountTooBig   = errors.New("discount cannot be higher minBasketValue")
	ErrCouponCodeAlreadyExist = errors.New("coupon code already used for another coupon")
)

// ErrApplyDiscountLessMin Custom error to apply function
type ErrApplyDiscountLessMin struct {
	MinValue int
	Current  int
}

// Error return the message with the value
func (e *ErrApplyDiscountLessMin) Error() string {
	return fmt.Sprintf(
		"cannot apply discount: value %d did not reach the minimum %d",
		e.Current, e.MinValue,
	)
}
