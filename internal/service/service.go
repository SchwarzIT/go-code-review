package service

import (
	"coupon_service/internal/repository/memdb"
	"fmt"

	"github.com/google/uuid"
)

// Repository interface to memdb repository
type Repository interface {
	FindByCode(string) (*memdb.Coupon, error)
	Save(*memdb.Coupon) error
}

// Service manage application features
type Service struct {
	repo Repository
}

// New create a new Service
func New(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// ApplyCoupon in the basket provided
// It returns an error if code coupon not exist and basket value must be positive value
func (s *Service) ApplyCoupon(basket *Basket, code string) error {
	coupon, err := s.repo.FindByCode(code)
	if err != nil {
		return err
	}

	if basket.Value <= 0 {
		return ErrApplyDiscount
	}

	if basket.Value < coupon.MinBasketValue {
		return &ErrApplyDiscountLessMin{
			MinValue: coupon.MinBasketValue,
			Current:  basket.Value,
		}
	}

	basket.AppliedDiscount = coupon.Discount
	basket.ApplicationSuccessful = true

	return nil
}

// CreateCoupon a new coupon
// It returns an error if discount not be a positive number and minBasketValue cant not be negative
// It returns an error if discount be higher that min basket
func (s *Service) CreateCoupon(discount int, code string, minBasketValue int) (*memdb.Coupon, error) {
	if discount <= 0 {
		return nil, ErrCouponDiscountValue
	}

	if minBasketValue < 0 {
		return nil, ErrCouponMinBasketValue
	}

	if discount > minBasketValue {
		return nil, ErrCouponDiscountTooBig
	}

	if _, err := s.repo.FindByCode(code); err == nil {
		return nil, ErrCouponCodeAlreadyExist
	}

	coupon := &memdb.Coupon{
		ID:             uuid.NewString(),
		Code:           code,
		Discount:       discount,
		MinBasketValue: minBasketValue,
	}

	if err := s.repo.Save(coupon); err != nil {
		return nil, fmt.Errorf("failed to save coupon: %w", err)
	}
	return coupon, nil
}

// GetCoupons return a list of coupons based on the codes provided
// It returns an error if case one of the code does not exist will
func (s *Service) GetCoupons(codes []string) ([]*memdb.Coupon, error) {
	coupons := make([]*memdb.Coupon, 0, len(codes))
	var errs []error

	for _, code := range codes {
		coupon, err := s.repo.FindByCode(code)
		if err != nil {
			errs = append(errs, fmt.Errorf("coupon code: %s, err: %s", code, err.Error()))
			continue
		}
		coupons = append(coupons, coupon)
	}

	if len(errs) > 0 {
		return coupons, fmt.Errorf("one or more errors occurred: %v", errs)
	}

	return coupons, nil
}
