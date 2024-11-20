package memdb

import (
	"fmt"

	"coupon_service/internal/entity"
)

type Config struct{}

type Repository struct {
	entries map[string]entity.Coupon
}

func New() *Repository {
	return &Repository{}
}

func (r *Repository) FindByCode(code string) (*entity.Coupon, error) {
	coupon, ok := r.entries[code]
	if !ok {
		return nil, fmt.Errorf("Coupon not found")
	}
	return &coupon, nil
}

func (r *Repository) Save(coupon *entity.Coupon) error {
	if r.entries == nil {
		r.entries = make(map[string]entity.Coupon)
	}

	if existingCoupon, exists := r.entries[coupon.Code]; exists {
		return fmt.Errorf("coupon with code %s already exists: %v", coupon.Code, existingCoupon)
	}

	r.entries[coupon.Code] = *coupon
	return nil
}
