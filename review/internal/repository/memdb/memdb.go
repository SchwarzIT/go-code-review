package memdb

import (
	"context"
	"fmt"
	"sync"

	"coupon_service/internal/service/entity"
)

type Repository struct {
	sync.RWMutex
	entries map[string]entity.Coupon
}

func New() *Repository {
	return &Repository{
		entries: make(map[string]entity.Coupon),
	}
}

func (r *Repository) FindByCode(ctx context.Context, code string) (entity.Coupon, error) {
	r.RLock()
	defer r.RUnlock()
	coupon, ok := r.entries[code]
	if !ok {
		return entity.Coupon{}, fmt.Errorf("coupon not found")
	}
	return coupon, nil
}

func (r *Repository) List(ctx context.Context, filter ...string) ([]entity.Coupon, error) {
	r.RLock()
	defer r.RUnlock()

	// If filters are provided, return only the matching coupons
	if len(filter) > 0 {
		coupons := make([]entity.Coupon, 0, len(filter))
		for _, code := range filter {
			if coupon, ok := r.entries[code]; ok {
				coupons = append(coupons, coupon)
			}
		}
		return coupons, nil
	}

	// If no filters are provided, return all coupons
	coupons := make([]entity.Coupon, 0, len(r.entries))
	for _, coupon := range r.entries {
		coupons = append(coupons, coupon)
	}

	return coupons, nil
}

func (r *Repository) Save(ctx context.Context, coupon entity.Coupon) error {
	r.Lock()
	defer r.Unlock()
	if _, ok := r.entries[coupon.Code]; ok {
		return fmt.Errorf("coupon with provided code already exists")
	}
	r.entries[coupon.Code] = coupon
	return nil
}

func (r *Repository) Delete(ctx context.Context, code string) error {
	r.Lock()
	defer r.Unlock()
	delete(r.entries, code)
	return nil
}
