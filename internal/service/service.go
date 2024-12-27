package service

import (
	. "coupon_service/internal/service/entity"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type Repository interface {
	FindByCode(string) (*Coupon, error)
	Save(*Coupon) error
}

type Service struct {
	repo Repository
}

func New(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// Custom errors for better error handling.
var (
	ErrApllyDiscount        = errors.New("cannot apply discount to a basket of non-positive value")
	ErrCouponDiscountValue  = errors.New("discount must be a positive integer")
	ErrCouponMinBasketValue = errors.New("minBasketValue cannot be negative")
)

func (s Service) ApplyCoupon(basket *Basket, code string) error {
	coupon, err := s.repo.FindByCode(code)
	if err != nil {
		return err
	}

	if basket.Value <= 0 {
		return ErrApllyDiscount
	}

	basket.AppliedDiscount = coupon.Discount
	basket.ApplicationSuccessful = true

	return nil
}

func (s *Service) CreateCoupon(discount int, code string, minBasketValue int) (*Coupon, error) {
	if discount <= 0 {
		return nil, ErrCouponDiscountValue
	}

	if minBasketValue < 0 {
		return nil, ErrCouponMinBasketValue
	}

	coupon := &Coupon{
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

func (s Service) GetCoupons(codes []string) ([]*Coupon, error) {
	coupons := make([]*Coupon, 0, len(codes))
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
