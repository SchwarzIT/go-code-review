package memdb

import (
	"errors"
	"fmt"
	"sync"

	"coupon_service/internal/service/entity"
)

// Repository defines the in-memory storage for Coupons.
// It implements the repository interface.
type Repository struct {
	entries map[string]*entity.Coupon
	mu      sync.RWMutex
}

// NewRepository creates and returns a new Repository instance.
func NewRepository() *Repository {
	return &Repository{
		entries: make(map[string]*entity.Coupon),
	}
}

// RepositoryInterface defines the methods that the Repository implements.
// Exported for external usage if needed.
type RepositoryInterface interface {
	FindByCode(string) (*entity.Coupon, error)
	Save(*entity.Coupon) error
	Delete(string) error
}

// Custom errors for better error handling.
var (
	ErrCouponNotFound = errors.New("coupon not found")
	ErrInvalidCoupon  = errors.New("invalid coupon")
)

// FindByCode retrieves a Coupon by its code.
// It returns a copy of the Coupon to prevent external modifications.
func (r *Repository) FindByCode(code string) (*entity.Coupon, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	coupon, exists := r.entries[code]
	if !exists {
		return nil, ErrCouponNotFound
	}

	// Return a copy to maintain immutability.
	couponCopy := *coupon
	return &couponCopy, nil
}

// Save stores a Coupon in the repository.
// It returns an error if the coupon is nil or has an empty code.
func (r *Repository) Save(coupon *entity.Coupon) error {
	if coupon == nil {
		return ErrInvalidCoupon
	}
	if coupon.Code == "" {
		return fmt.Errorf("%w: coupon code is empty", ErrInvalidCoupon)
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	// Store a copy to prevent external modifications affecting the repository.
	couponCopy := *coupon
	r.entries[coupon.Code] = &couponCopy
	return nil
}

// Delete removes a Coupon from the repository by its code.
// It returns an error if the coupon does not exist.
func (r *Repository) Delete(code string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.entries[code]; !exists {
		return ErrCouponNotFound
	}

	delete(r.entries, code)
	return nil
}
