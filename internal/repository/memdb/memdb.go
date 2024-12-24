package memdb

import (
	"coupon_service/internal/service/entity"
	"fmt"
)

type Config struct{}

type repository interface {
	FindByCode(string) (*entity.Coupon, error)
	Save(entity.Coupon) error
}

type Repository struct {
	entries map[string]entity.Coupon
}

func New() *Repository {
	return &Repository{
		entries: make(map[string]entity.Coupon),
	}
}

func (r *Repository) FindByCode(code string) (*entity.Coupon, error) {
	coupon, ok := r.entries[code]
	if !ok {
		return nil, fmt.Errorf("Coupon not found")
	}
	return &coupon, nil
}

func (r *Repository) Save(coupon entity.Coupon) error {
	r.entries[coupon.Code] = coupon
	return nil
}
