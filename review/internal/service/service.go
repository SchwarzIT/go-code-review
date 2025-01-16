package service

import (
	"fmt"

	"github.com/google/uuid"

	"coupon_service/internal/service/entity"
)

type Repository interface {
	FindByCode(string) (entity.Coupon, error)
	List(...string) ([]entity.Coupon, error)
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

func (s *Service) ApplyCoupon(value, appliedDiscount int, code string) (entity.Basket, error) {
	basket := entity.Basket{
		Value:           value,
		AppliedDiscount: appliedDiscount,
	}

	// return basket without changes if it's empty
	if basket.Value <= 0 {
		return basket, nil
	}

	coupon, err := s.repo.FindByCode(code)
	if err != nil {
		return entity.Basket{}, fmt.Errorf("failed to get coupon: %w", err)
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
		return entity.Basket{}, fmt.Errorf("failed to delete coupon: %w", err)
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

func (s *Service) ListCoupons(codes ...string) ([]entity.Coupon, error) {
	coupons, err := s.repo.List(codes...)
	if err != nil {
		return nil, fmt.Errorf("failed to get coupons: %w", err)
	}

	return coupons, nil
}
