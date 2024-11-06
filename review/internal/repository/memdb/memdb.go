package memdb

import (
	"coupon_service/internal/service/entity"
	"errors"
	"fmt"
	"sync"
)

// Custom error types
var (
	ErrCouponNotFound    = errors.New("coupon not found")
	ErrInvalidCouponType = errors.New("invalid coupon type")
	ErrDuplicateCoupon   = errors.New("coupon code already exists")
)

// Repository is an in-memory repository for coupons
type Repository struct {
	entries sync.Map
}

// New creates a new Repository
func New() *Repository {
	return &Repository{}
}

// FindByCode finds a coupon by its code
func (r *Repository) FindByCode(code string) (*entity.Coupon, error) {
	value, ok := r.entries.Load(code)
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrCouponNotFound, code)
	}
	coupon, ok := value.(entity.Coupon)
	if !ok {
		return nil, fmt.Errorf("%w: %v", ErrInvalidCouponType, value)
	}
	return &coupon, nil
}

// Save saves a coupon
func (r *Repository) Save(coupon entity.Coupon) error {
	_, loaded := r.entries.LoadOrStore(coupon.Code, coupon)
	if loaded {
		return fmt.Errorf("%w: %s", ErrDuplicateCoupon, coupon.Code)
	}
	return nil
}
