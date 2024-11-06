package service

import (
	"coupon_service/internal/repository/memdb"
	"coupon_service/internal/service/entity"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		config Config
	}
	tests := []struct {
		name string
		args args
		want *Service
	}{
		{"initialize service with nil repo", args{config: Config{CouponsRepository: nil}}, &Service{repo: nil}},
		{"initialize service with memdb repo", args{config: Config{CouponsRepository: memdb.New()}}, &Service{repo: memdb.New()}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.config); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_ApplyCoupon(t *testing.T) {
	repo := memdb.New()
	s := New(Config{CouponsRepository: repo})

	// Create a coupon to be applied
	coupon := entity.Coupon{
		Code:           "DISCOUNT10",
		Discount:       10,
		MinBasketValue: 50,
	}
	repo.Save(coupon)

	type args struct {
		basket entity.Basket
		code   string
	}
	tests := []struct {
		name    string
		args    args
		wantB   *entity.Basket
		wantErr bool
	}{
		{
			name: "apply valid coupon",
			args: args{
				basket: entity.Basket{Value: 100},
				code:   "DISCOUNT10",
			},
			wantB:   &entity.Basket{Value: 100, AppliedDiscount: 10, ApplicationSuccessful: true},
			wantErr: false,
		},
		{
			name: "apply invalid coupon",
			args: args{
				basket: entity.Basket{Value: 100},
				code:   "INVALID",
			},
			wantB:   nil,
			wantErr: true,
		},
		{
			name: "apply coupon to basket with insufficient value",
			args: args{
				basket: entity.Basket{Value: 0},
				code:   "DISCOUNT10",
			},
			wantB:   nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotB, err := s.ApplyCoupon(tt.args.basket, tt.args.code)
			if (err != nil) != tt.wantErr {
				t.Errorf("ApplyCoupon() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotB, tt.wantB) {
				t.Errorf("ApplyCoupon() gotB = %v, want %v", gotB, tt.wantB)
			}
		})
	}
}

func TestService_CreateCoupon(t *testing.T) {
	repo := memdb.New()
	s := New(Config{CouponsRepository: repo})

	type args struct {
		discount       int
		code           string
		minBasketValue int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"create valid coupon", args{10, "DISCOUNT10", 50}, false},
		{"create duplicate coupon", args{10, "DISCOUNT10", 50}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := s.CreateCoupon(tt.args.discount, tt.args.code, tt.args.minBasketValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateCoupon() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestService_GetCoupons(t *testing.T) {
	repo := memdb.New()
	s := New(Config{CouponsRepository: repo})

	// Create some coupons
	coupon1 := entity.Coupon{
		Code:           "DISCOUNT10",
		Discount:       10,
		MinBasketValue: 50,
	}
	coupon2 := entity.Coupon{
		Code:           "DISCOUNT20",
		Discount:       20,
		MinBasketValue: 100,
	}
	repo.Save(coupon1)
	repo.Save(coupon2)

	type args struct {
		codes []string
	}
	tests := []struct {
		name    string
		args    args
		want    []entity.Coupon
		wantErr bool
	}{
		{
			name:    "get existing coupons",
			args:    args{codes: []string{"DISCOUNT10", "DISCOUNT20"}},
			want:    []entity.Coupon{coupon1, coupon2},
			wantErr: false,
		},
		{
			name:    "get non-existing coupon",
			args:    args{codes: []string{"INVALID"}},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "get mixed existing and non-existing coupons",
			args:    args{codes: []string{"DISCOUNT10", "INVALID"}},
			want:    []entity.Coupon{coupon1},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.GetCoupons(tt.args.codes)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCoupons() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCoupons() got = %v, want %v", got, tt.want)
			}
		})
	}
}
