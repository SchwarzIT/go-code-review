package service

import (
	. "coupon_service/internal/service/entity"
	"fmt"

	"github.com/google/uuid"
)

type couponsRepository interface {
	FindByCode(code string) (*Coupon, error)
	Save(coupon Coupon) error
}

type Config struct {
	CouponsRepository couponsRepository
}

type Service struct {
	repo couponsRepository
}

func New(config Config) *Service {
	return &Service{
		repo: config.CouponsRepository,
	}
}

func (s *Service) ApplyCoupon(basket Basket, code string) (*Basket, error) {
	coupon, err := s.repo.FindByCode(code)
	if err != nil {
		return nil, err
	}

	if basket.Value > 0 {
		basket.AppliedDiscount = coupon.Discount
		basket.ApplicationSuccessful = true
		return &basket, nil
	}

	return nil, fmt.Errorf("Tried to apply discount to negative value")
}

func (s *Service) CreateCoupon(discount int, code string, minBasketValue int) error {
	coupon := Coupon{
		Discount:       discount,
		Code:           code,
		MinBasketValue: minBasketValue,
		ID:             uuid.NewString(),
	}

	if err := s.repo.Save(coupon); err != nil {
		return err
	}
	return nil
}

func (s *Service) GetCoupons(codes []string) ([]Coupon, error) {
	coupons := make([]Coupon, 0, len(codes))
	var errors []string

	for idx, code := range codes {
		coupon, err := s.repo.FindByCode(code)
		if err != nil {
			errors = append(errors, fmt.Sprintf("code: %s, index: %d", code, idx))
			continue
		}
		coupons = append(coupons, *coupon)
	}

	if len(coupons) == 0 {
		return nil, fmt.Errorf("errors: %s", errors)
	}

	return coupons, nil
}
