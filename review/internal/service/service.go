package service

import (
	"fmt"

	"github.com/google/uuid"

	"coupon_service/internal/service/entity"
)

type Repository interface {
	FindByCode(string) (entity.Coupon, error)
	Save(entity.Coupon) error
	Delete(string) error
}

type Service struct {
	repo Repository
}

func New(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) ApplyCoupon(basket entity.Basket, code string) (entity.Basket, error) {
	// return basket without changes if it's empty
	if basket.Value <= 0 {
		return basket, nil
	}

	coupon, err := s.repo.FindByCode(code)
	if err != nil {
		// TODO: define code type and implement Stringer interface to avoid printing full code in logs
		return entity.Basket{}, fmt.Errorf("failed to get coupon by %q code: %w", code, err)
	}

	// check if we fit MinBasketValue constraint
	if coupon.MinBasketValue > basket.Value {
		return entity.Basket{}, fmt.Errorf("not enough value in basket: should be gte %d", coupon.MinBasketValue)
	}

	diff := basket.Value - coupon.Discount
	if diff < 0 {
		basket.Value = 0
		// apply discount until 0 value, the rest of the points are burned out
		basket.AppliedDiscount += coupon.Discount + diff
	} else {
		basket.Value = diff
		// apply full coupon
		basket.AppliedDiscount += coupon.Discount
	}
	basket.ApplicationSuccessful = true

	// delete coupon because of successfull apply
	if err := s.repo.Delete(code); err != nil {
		// TODO: define code type and implement Stringer interface to avoid printing full code in logs
		return entity.Basket{}, fmt.Errorf("failed to delete coupon by %q code: %w", code, err)
	}

	return basket, nil
}

func (s *Service) CreateCoupon(discount int, code string, minBasketValue int) error {
	id, err := uuid.NewRandom()
	if err != nil {
		return fmt.Errorf("failed to generate coupon id: %w", err)
	}

	coupon := entity.Coupon{
		Discount:       discount,
		Code:           code,
		MinBasketValue: minBasketValue,
		ID:             id.String(),
	}

	if err := s.repo.Save(coupon); err != nil {
		return fmt.Errorf("failed to get coupon: %w", err)
	}
	return nil
}

func (s *Service) GetCoupons(codes []string) ([]entity.Coupon, error) {
	coupons := make([]entity.Coupon, 0, len(codes))

	for _, code := range codes {
		coupon, err := s.repo.FindByCode(code)
		if err != nil {
			// TODO: define code type and implement Stringer interface to avoid printing full code in logs
			return nil, fmt.Errorf("failed to get coupon by %q code: %w", code, err)
		}
		coupons = append(coupons, coupon)
	}

	return coupons, nil
}
