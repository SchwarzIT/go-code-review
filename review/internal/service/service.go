package service

import (
	"fmt"
	"regexp"

	"github.com/google/uuid"

	"coupon_service/internal/entity"
)

type Repository interface {
	FindByCode(string) (*entity.Coupon, error)
	Save(*entity.Coupon) error
}

type Service struct {
	repo Repository
}

// New creates a new Service instance with the provided Repository.
func New(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// ApplyCoupon attempts to apply a discount coupon to a given basket.
// If the coupon is found and meets requirements, it updates the basket's discount.
func (s *Service) ApplyCoupon(basket *entity.Basket, code string) error {
	if !validateCode(code) {
		return fmt.Errorf("invalid code")
	}

	coupon, err := s.repo.FindByCode(code)
	if err != nil {
		return fmt.Errorf("error fetching coupon: %w", err)
	}

	if coupon == nil {
		return fmt.Errorf("coupon with code: %s not found", code)
	}

	if basket.Value < coupon.MinBasketValue {
		return fmt.Errorf("basket value below minimum required for coupon")
	}

	basket.AppliedDiscount = coupon.Discount
	basket.ApplicationSuccessful = true

	return nil
}

// CreateCoupon saves a new coupon with the specified discount, code, and minimum basket value.
// Returns an error if the coupon could not be saved.
func (s *Service) CreateCoupon(discount int, code string, minBasketValue int) error {
	if discount < 0 {
		return fmt.Errorf("discount cannot be negative")
	}

	if minBasketValue < 0 {
		return fmt.Errorf("minBasketValue cannot be negative")
	}

	if !validateCode(code) {
		return fmt.Errorf("invalid code")
	}

	coupon := &entity.Coupon{
		Discount:       discount,
		Code:           code,
		MinBasketValue: minBasketValue,
		ID:             uuid.NewString(),
	}

	if err := s.repo.Save(coupon); err != nil {
		return fmt.Errorf("failed to save coupon: %w", err)
	}

	return nil
}

// GetCoupons retrieves a list of coupons for the given codes, with a separate list of codes not found.
func (s *Service) GetCoupons(codes []string) []entity.Coupon {
	var coupons []entity.Coupon

	for _, code := range codes {
		coupon, err := s.repo.FindByCode(code)
		if err == nil && coupon != nil {
			coupons = append(coupons, *coupon)
		}
	}

	return coupons
}

// GetCoupon retrieves a single coupon by its code. Returns an error if the coupon is not found.
func (s *Service) GetCoupon(code string) (*entity.Coupon, error) {
	coupon, err := s.repo.FindByCode(code)
	if err != nil {
		return nil, fmt.Errorf("error retrieving coupon: %w", err)
	}
	if coupon == nil {
		return nil, fmt.Errorf("coupon with code %s not found", code)
	}
	return coupon, err
}

// validateCode checks if the given code meets certain criteria.
func validateCode(code string) bool {
	if len(code) < 5 {
		return false
	}

	match, err := regexp.MatchString(`^[a-zA-Z0-9]+$`, code)
	if err != nil {
		return false
	}

	return match
}
