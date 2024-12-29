package memdb

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepository_Save(t *testing.T) {
	repo := SetupTestRepository(t)

	t.Run("Save and retrieve a valid coupon", func(t *testing.T) {
		coupon := Coupon{
			ID:             "1",
			Code:           "DISCOUNT10",
			Discount:       10,
			MinBasketValue: 50,
		}

		err := repo.Save(&coupon)
		assert.NoError(t, err, "Expected no error when saving a valid coupon")

		storedCoupon, err := repo.FindByCode("DISCOUNT10")
		assert.NoError(t, err, "Expected to find the saved coupon")
		assert.Equal(t, coupon, storedCoupon, "Stored coupon should match the saved coupon")
	})

	t.Run("Overwrite an existing coupon", func(t *testing.T) {
		coupon1 := Coupon{
			ID:             "2",
			Code:           "SAVE20",
			Discount:       20,
			MinBasketValue: 100,
		}

		coupon2 := Coupon{
			ID:             "3",
			Code:           "SAVE20",
			Discount:       25,
			MinBasketValue: 150,
		}

		err := repo.Save(&coupon1)
		assert.NoError(t, err, "Expected no error when saving the first coupon")

		err = repo.Save(&coupon2)
		assert.NoError(t, err, "Expected no error when overwriting the existing coupon")

		storedCoupon, err := repo.FindByCode("SAVE20")
		assert.NoError(t, err, "Expected to find the overwritten coupon")
		assert.Equal(t, coupon2, storedCoupon, "Stored coupon should match the latest saved coupon")
	})

	t.Run("Save a coupon with empty code should fail", func(t *testing.T) {
		coupon := Coupon{
			ID:             "4",
			Code:           "",
			Discount:       15,
			MinBasketValue: 75,
		}

		err := repo.Save(&coupon)
		assert.Error(t, err, "Expected an error when saving a coupon with empty code")
		assert.Equal(t, ErrInvalidCoupon.Error()+": coupon code is empty", err.Error(), "Error message should indicate empty coupon code")
	})

	t.Run("Save a nil coupon should fail", func(t *testing.T) {
		err := repo.Save(nil)
		assert.Error(t, err, "Expected an error when saving a nil coupon")
		assert.Equal(t, ErrInvalidCoupon, err, "Error should be ErrInvalidCoupon")
	})
}

func TestRepository_FindByCode(t *testing.T) {
	repo := SetupTestRepository(t)

	t.Run("Find existing coupon", func(t *testing.T) {
		coupon := Coupon{
			ID:             "5",
			Code:           "SUMMER20",
			Discount:       20,
			MinBasketValue: 100,
		}

		err := repo.Save(&coupon)
		assert.NoError(t, err, "Expected no error when saving a valid coupon")

		retrievedCoupon, err := repo.FindByCode("SUMMER20")
		assert.NoError(t, err, "Expected to find the existing coupon")
		assert.Equal(t, coupon, retrievedCoupon, "Retrieved coupon should match the saved coupon")
	})

	t.Run("Attempt to find a non-existent coupon", func(t *testing.T) {
		_, err := repo.FindByCode("NONEXISTENT")
		assert.Error(t, err, "Expected an error when finding a non-existent coupon")
		assert.Equal(t, ErrCouponNotFound, err, "Error should be ErrCouponNotFound")
	})
}

func TestRepository_Concurrency(t *testing.T) {
	repo := SetupTestRepository(t)
	var wg sync.WaitGroup
	numGoroutines := 50
	couponCodePrefix := "CONCUR_"

	saveCoupons := func(start, end int) {
		defer wg.Done()
		for i := start; i < end; i++ {
			coupon := &Coupon{
				ID:             fmt.Sprintf("%d", i),
				Code:           fmt.Sprintf("%s%d", couponCodePrefix, i),
				Discount:       i % 100,
				MinBasketValue: i * 10,
			}
			err := repo.Save(coupon)
			if err != nil {
				// It's generally not recommended to use assert within goroutines.
				// Instead, you might want to log the error or handle it accordingly.
				// For simplicity, we'll use t.Errorf here, but be aware of potential race conditions.
				t.Errorf("Failed to save coupon %s: %v", coupon.Code, err)
			}
		}
	}

	findCoupons := func(start, end int) {
		defer wg.Done()
		for i := start; i < end; i++ {
			_, _ = repo.FindByCode(fmt.Sprintf("%s%d", couponCodePrefix, i))
			// Ignoring errors as some coupons might not have been saved yet.
		}
	}

	// Launch goroutines to save coupons
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go saveCoupons(i*100, (i+1)*100)
	}

	// Launch goroutines to find coupons
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go findCoupons(i*100, (i+1)*100)
	}

	wg.Wait()

	// Verify a few coupons to ensure they were saved correctly
	for i := 0; i < numGoroutines*100; i += 1000 {
		code := fmt.Sprintf("%s%d", couponCodePrefix, i)
		coupon, err := repo.FindByCode(code)
		if assert.NoError(t, err, "Expected to find the coupon %s", code) {
			assert.Equal(t, code, coupon.Code, "Coupon code should match")
			assert.Equal(t, fmt.Sprintf("%d", i), coupon.ID, "Coupon ID should match")
			assert.Equal(t, i%100, coupon.Discount, "Coupon discount should match")
			assert.Equal(t, i*10, coupon.MinBasketValue, "Coupon min basket value should match")
		}
	}
}
