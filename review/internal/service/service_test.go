package service

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/mock"

	"coupon_service/internal/mocks"
	"coupon_service/internal/service/entity"
)

func TestService_ApplyCoupon(t *testing.T) {
	mockRepo := mocks.NewMockRepository(t)

	mockRepo.On("FindByCode", "DISCOUNT10").Return(entity.Coupon{
		Code:           "DISCOUNT10",
		Discount:       10,
		MinBasketValue: 50,
	}, nil)
	mockRepo.On("FindByCode", "INVALID").Return(entity.Coupon{}, fmt.Errorf("coupon not found"))

	mockRepo.On("Delete", mock.Anything).Return(nil)

	svc := New(mockRepo)

	tests := []struct {
		name       string
		value      int
		discount   int
		code       string
		wantBasket entity.Basket
		wantErr    bool
	}{
		{
			name:     "successful coupon application",
			value:    100,
			discount: 10,
			code:     "DISCOUNT10",
			wantBasket: entity.Basket{
				Value:                 90,
				AppliedDiscount:       20,
				ApplicationSuccessful: true,
			},
			wantErr: false,
		},
		{
			name:       "coupon not found",
			value:      100,
			discount:   10,
			code:       "INVALID",
			wantBasket: entity.Basket{},
			wantErr:    true,
		},
		{
			name:       "basket value less than minimum",
			value:      30,
			discount:   10,
			code:       "DISCOUNT10",
			wantBasket: entity.Basket{},
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBasket, err := svc.ApplyCoupon(tt.value, tt.discount, tt.code)
			if (err != nil) != tt.wantErr {
				t.Errorf("ApplyCoupon() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotBasket, tt.wantBasket) {
				t.Errorf("ApplyCoupon() = %v, want %v", gotBasket, tt.wantBasket)
			}
		})
	}
}

func TestService_CreateCoupon(t *testing.T) {
	mockRepo := mocks.NewMockRepository(t)
	mockRepo.On("Save", mock.Anything).Return(nil).Once()
	mockRepo.On("Save", mock.MatchedBy(func(c entity.Coupon) bool {
		return c.Code == "DISCOUNT15"
	})).Return(errors.New("coupon already exists")).Once()

	svc := New(mockRepo)

	tests := []struct {
		name           string
		discount       int
		code           string
		minBasketValue int
		wantErr        bool
	}{
		{
			name:           "successful coupon creation",
			discount:       15,
			code:           "DISCOUNT15",
			minBasketValue: 100,
			wantErr:        false,
		},
		{
			name:           "duplicate coupon code",
			discount:       10,
			code:           "DISCOUNT15",
			minBasketValue: 50,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := svc.CreateCoupon(tt.discount, tt.code, tt.minBasketValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateCoupon() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
