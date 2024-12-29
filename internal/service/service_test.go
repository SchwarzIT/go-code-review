package service

import (
	"coupon_service/internal/repository/memdb"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestService_CreateCoupon(t *testing.T) {
	repo := memdb.SetupTestRepository(t)
	sev := New(repo)
	tests := []struct {
		name  string
		input struct {
			discount       int
			code           string
			minBasketValue int
		}
		expectedError error
	}{
		{
			name: "Valid coupon",
			input: struct {
				discount       int
				code           string
				minBasketValue int
			}{
				discount:       20,
				code:           "AA1",
				minBasketValue: 30,
			},
			expectedError: nil,
		},
		{
			name: "Invalid coupon discount value is zero",
			input: struct {
				discount       int
				code           string
				minBasketValue int
			}{
				discount:       0,
				code:           "AA2",
				minBasketValue: 30,
			},
			expectedError: ErrCouponDiscountValue,
		},
		{
			name: "Invalid coupon minBasketValue is negative",
			input: struct {
				discount       int
				code           string
				minBasketValue int
			}{
				discount:       10,
				code:           "AA3",
				minBasketValue: -1,
			},
			expectedError: ErrCouponMinBasketValue,
		},
		{
			name: "Invalid coupon discount is higher than minBasket",
			input: struct {
				discount       int
				code           string
				minBasketValue int
			}{
				discount:       50,
				code:           "AA4",
				minBasketValue: 30,
			},
			expectedError: ErrCouponDiscountTooBig,
		},
		{
			name: "Invalid coupon discount code already exist",
			input: struct {
				discount       int
				code           string
				minBasketValue int
			}{
				discount:       20,
				code:           "AA1",
				minBasketValue: 30,
			},
			expectedError: ErrCouponCodeAlreadyExist,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			coupon, err := sev.CreateCoupon(tt.input.discount, tt.input.code, tt.input.minBasketValue)

			if tt.expectedError != nil {
				assert.Error(t, err, "Expected an error")
				assert.Equal(t, tt.expectedError.Error(), err.Error(), fmt.Sprintf("Expected this error %s", tt.expectedError.Error()))

			} else {
				assert.NoError(t, err, "Did not expect an error for this input %v", tt.input)
				assert.NotNil(t, coupon, "Expect not nill value in coupon return")
			}
		})
	}

}

func TestService_GetCoupons(t *testing.T) {
	repo := memdb.SetupTestRepository(t)
	svc := New(repo)

	coupon, err := svc.CreateCoupon(20, "AA1", 100)
	assert.NoError(t, err, "Expected no error")
	assert.NotNil(t, coupon, "Expected no nil value")

	tests := []struct {
		name           string
		input          []string
		expectedOutput []string
		expectedError  bool
	}{
		{
			name:          "Get valid coupon",
			input:         []string{"AA1"},
			expectedError: false,
		},
		{
			name:          "Get valid coupon",
			input:         []string{"AA2"},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rCoupons, err := svc.GetCoupons(tt.input)
			checkCoupons := func(code string) bool {
				for _, iCoupon := range rCoupons {
					if iCoupon.Code == code {
						return true
					}
				}
				return false
			}
			if tt.expectedError {
				assert.Error(t, err, "Expected an error")
			} else {
				for _, item := range tt.input {
					assert.True(t, checkCoupons(item), fmt.Sprintf("Expected found this coupon code %s", item))
				}
			}
		})
	}
}

func TestService_ApplyCoupon(t *testing.T) {
	repo := memdb.SetupTestRepository(t)
	svc := New(repo)

	coupon, err := svc.CreateCoupon(20, "AA1", 50)
	assert.NoError(t, err, "Expected no error")
	assert.NotNil(t, coupon, "Expected no nil value")

	tests := []struct {
		name  string
		input struct {
			basket *Basket
			code   string
		}
		expectedDiscount int
		expectedError    error
	}{
		{
			name: "Apply valid discount",
			input: struct {
				basket *Basket
				code   string
			}{
				basket: &Basket{
					Value: 100,
				},
				code: "AA1",
			},
			expectedDiscount: 20,
			expectedError:    nil,
		},
		{
			name: "Apply invalid discount to not exist coupon",
			input: struct {
				basket *Basket
				code   string
			}{
				basket: &Basket{
					Value: 100,
				},
				code: "AA2",
			},
			expectedError: memdb.ErrCouponNotFound,
		},
		{
			name: "Apply invalid discount to basket value equal zero",
			input: struct {
				basket *Basket
				code   string
			}{
				basket: &Basket{
					Value: 0,
				},
				code: "AA1",
			},
			expectedError: ErrApplyDiscount,
		},
		{
			name: "Apply invalid discount to basket value less than minimum",
			input: struct {
				basket *Basket
				code   string
			}{
				basket: &Basket{
					Value: 20,
				},
				code: "AA1",
			},
			expectedError: &ErrApplyDiscountLessMin{
				Current:  20,
				MinValue: 50,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err = svc.ApplyCoupon(tt.input.basket, tt.input.code)
			if tt.expectedError != nil {
				assert.True(t, strings.Contains(err.Error(), tt.expectedError.Error()), fmt.Sprintf("Expected has the error %s", tt.expectedError.Error()))

			} else {
				assert.Equal(t, tt.expectedDiscount, tt.input.basket.AppliedDiscount, "Expected that the discount has been applied")
			}
		})
	}
}
