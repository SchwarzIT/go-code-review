package memdb

import (
	"coupon_service/internal/service/entity"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepository_Save(t *testing.T) {
	repo := New()

	// Test saving a valid coupon
	coupon := entity.Coupon{
		Code:           "DISCOUNT10",
		Discount:       10,
		MinBasketValue: 50,
	}
	err := repo.Save(coupon)
	assert.NoError(t, err, "expected no error when saving a valid coupon")

	// Test saving a duplicate coupon
	err = repo.Save(coupon)
	assert.ErrorIs(t, err, ErrDuplicateCoupon, "expected duplicate coupon error")
}

func TestRepository_FindByCode(t *testing.T) {
	repo := New()

	// Test finding a non-existent coupon
	_, err := repo.FindByCode("NONEXISTENT")
	assert.ErrorIs(t, err, ErrCouponNotFound, "expected coupon not found error")

	// Test finding a valid coupon
	coupon := entity.Coupon{
		Code:           "DISCOUNT10",
		Discount:       10,
		MinBasketValue: 50,
	}
	err = repo.Save(coupon)
	assert.NoError(t, err, "expected no error when saving a valid coupon")

	foundCoupon, err := repo.FindByCode("DISCOUNT10")
	assert.NoError(t, err, "expected no error when finding a valid coupon")
	assert.Equal(t, coupon, *foundCoupon, "expected found coupon to match saved coupon")

	// Test finding a coupon with invalid type (simulated by storing a different type)
	repo.entries.Store("INVALID", "invalid type")
	_, err = repo.FindByCode("INVALID")
	assert.ErrorIs(t, err, ErrInvalidCouponType, "expected invalid coupon type error")
}
