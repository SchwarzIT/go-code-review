package memdb

import (
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

func (r *Repository) FindByCode(code string) (entity.Coupon, error) {
	r.RLock()
	defer r.RUnlock()
	coupon, ok := r.entries[code]
	if !ok {
		return entity.Coupon{}, fmt.Errorf("coupon not found")
	}
	return coupon, nil
}

func (r *Repository) Save(coupon entity.Coupon) error {
	r.Lock()
	defer r.Unlock()
	if _, ok := r.entries[coupon.Code]; ok {
		return fmt.Errorf("coupon with provided code already exists")
	}
	r.entries[coupon.Code] = coupon
	return nil
}

func (r *Repository) Delete(code string) error {
	r.Lock()
	defer r.Unlock()
	delete(r.entries, code)
	return nil
}
