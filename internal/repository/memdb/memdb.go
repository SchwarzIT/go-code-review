package memdb

import (
	"errors"
	"fmt"
	"sync"
)

type Coupon struct {
	ID             string `json:"id"`
	Code           string `json:"code"`
	Discount       int    `json:"discount"`
	MinBasketValue int    `json:"min_basket_value"`
}

// Repository defines the in-memory storage for Coupons.
// It implements the repository interface.
type Repository struct {
	entries map[string]*Coupon
	mu      sync.RWMutex
}

// NewRepository creates and returns a new Repository instance.
func NewRepository() *Repository {
	return &Repository{
		entries: make(map[string]*Coupon),
	}
}

// RepositoryInterface defines the methods that the Repository implements.
// Exported for external usage if needed.
type RepositoryInterface interface {
	FindByCode(string) (*Coupon, error)
	Save(*Coupon) error
	Delete(string) error
}

// Custom errors for better error handling.
var (
	ErrCouponNotFound = errors.New("coupon not found")
	ErrInvalidCoupon  = errors.New("invalid coupon")
)

// FindByCode retrieves a Coupon by its code.
// It returns a copy of the Coupon to prevent external modifications.
func (r *Repository) FindByCode(code string) (Coupon, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	coupon, exists := r.entries[code]
	if !exists {
		return Coupon{}, ErrCouponNotFound
	}

	// Return a copy to maintain immutability.
	return *coupon, nil
}

// Save stores a Coupon in the repository.
// It returns an error if the coupon is nil or has an empty code.
func (r *Repository) Save(coupon *Coupon) error {
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
