package service

import (
	"coupon_service/internal/entity"
	memdb "coupon_service/internal/repository"
	"errors"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type MockRepository struct {
	coupons map[string]entity.Coupon
}

func (m *MockRepository) FindByCode(code string) (*entity.Coupon, error) {
	if coupon, exists := m.coupons[code]; exists {
		return &coupon, nil
	}
	return nil, errors.New("coupon not found")
}

func (m *MockRepository) Save(coupon *entity.Coupon) error {
	if coupon == nil {
		return fmt.Errorf("nil coupon")
	}
	m.coupons[coupon.Code] = *coupon
	return nil
}

func TestNew(t *testing.T) {
	validRepo := memdb.New()
	tests := []struct {
		name string
		repo Repository
		want *Service
	}{
		{
			name: "Initialize service with nil repo",
			repo: Repository(nil),
			want: &Service{repo: Repository(nil)},
		},
		{
			name: "Initialize service with valid repo",
			repo: validRepo,
			want: &Service{repo: validRepo},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(tt.repo)
			assert.Equal(t, tt.want, got, "Test %q failed: expected service to be %v, got %v", tt.name, tt.want, got)
		})
	}
}

func TestService_ApplyCoupon(t *testing.T) {
	tests := []struct {
		name         string
		basket       entity.Basket
		coupon       entity.Coupon
		repo         Repository
		wantErr      bool
		wantDiscount int
	}{
		{
			name:   "Apply valid coupon",
			basket: entity.Basket{Value: 100},
			coupon: entity.Coupon{
				ID:             uuid.NewString(),
				Code:           "BlackFriday2024",
				Discount:       10,
				MinBasketValue: 30,
			},
			wantErr:      false,
			wantDiscount: 10,
		},
		{
			name:   "Invalid coupon code",
			basket: entity.Basket{Value: 100},
			coupon: entity.Coupon{
				ID:             uuid.NewString(),
				Code:           "",
				Discount:       10,
				MinBasketValue: 30,
			},
			wantErr:      true,
			wantDiscount: 10,
		},
		{
			name:   "Basket value below minimum",
			basket: entity.Basket{Value: 20},
			coupon: entity.Coupon{
				ID:             uuid.NewString(),
				Code:           "Discount20",
				Discount:       5,
				MinBasketValue: 30,
			},
			wantErr:      true,
			wantDiscount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockRepository{
				coupons: map[string]entity.Coupon{
					tt.coupon.Code: tt.coupon,
				},
			}
			s := Service{
				repo: mockRepo,
			}
			err := s.ApplyCoupon(&tt.basket, tt.coupon.Code)
			if tt.wantErr {
				assert.Error(t, err, "Expected an error but got none")
			} else {
				assert.NoError(t, err, "Expected no error but got one")
			}
		})
	}
}

func TestService_CreateCoupon(t *testing.T) {
	type args struct {
		discount       int
		code           string
		minBasketValue int
	}
	tests := []struct {
		name        string
		repo        Repository
		args        args
		wantErr     bool
		prepareRepo func(repo Repository)
	}{
		{
			name:    "Valid coupon creation",
			repo:    memdb.New(),
			args:    args{discount: 10, code: "Superdiscount", minBasketValue: 55},
			wantErr: false,
		},
		{
			name:    "Invalid discount (negative)",
			repo:    memdb.New(),
			args:    args{discount: -10, code: "InvalidDiscount", minBasketValue: 55},
			wantErr: true,
		},
		{
			name:    "Zero basket value (valid coupon)",
			repo:    memdb.New(),
			args:    args{discount: 0, code: "NoDiscount", minBasketValue: 0},
			wantErr: false,
		},
		{
			name:    "Duplicate coupon code",
			repo:    memdb.New(),
			args:    args{discount: 10, code: "Superdiscount", minBasketValue: 55},
			wantErr: true,
			prepareRepo: func(repo Repository) { // This is needed to have a duplicated coupon
				_ = repo.Save(&entity.Coupon{
					Code: "Superdiscount", Discount: 10, MinBasketValue: 55,
				})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.prepareRepo != nil {
				tt.prepareRepo(tt.repo)
			}

			s := Service{repo: tt.repo}
			err := s.CreateCoupon(tt.args.discount, tt.args.code, tt.args.minBasketValue)
			if tt.wantErr {
				assert.Error(t, err, "Expected an error but got none")
			} else {
				assert.NoError(t, err, "Expected no error but got one")
			}
		})
	}
}

func TestService_GetCoupons(t *testing.T) {
	tests := []struct {
		name        string
		codes       []string
		wantCoupons int
	}{
		{name: "Valid coupon", codes: []string{"BlackFriday2024"}, wantCoupons: 1},
		{name: "Multiple valid coupons", codes: []string{"BlackFriday2024", "StackIT"}, wantCoupons: 2},
		{name: "Invalid coupon", codes: []string{"InvalidCode"}, wantCoupons: 0},
		{name: "Multiple invalid coupons", codes: []string{"InvalidCode", ""}, wantCoupons: 0},
		{name: "Valid and invalid coupons", codes: []string{"BlackFriday2024", "InvalidCode"}, wantCoupons: 1},
	}

	mockRepo := &MockRepository{
		coupons: map[string]entity.Coupon{
			"BlackFriday2024": {
				ID:             "00001",
				Code:           "BlackFriday2024",
				Discount:       10,
				MinBasketValue: 30,
			},
			"StackIT": {
				ID:             "00002",
				Code:           "BlackFriday2024",
				Discount:       100,
				MinBasketValue: 150,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Service{repo: mockRepo}
			coupons := s.GetCoupons(tt.codes)

			assert.Equal(t, tt.wantCoupons, len(coupons))
		})

	}
}

func TestService_GetCoupon(t *testing.T) {
	tests := []struct {
		name           string
		code           string
		wantErr        bool
		expectedCoupon *entity.Coupon
	}{
		{name: "Valid coupon", code: "BlackFriday2024", wantErr: false, expectedCoupon: &entity.Coupon{
			ID:             "00001",
			Code:           "BlackFriday2024",
			Discount:       10,
			MinBasketValue: 30,
		}},
		{name: "Invalid coupon", code: "InvalidCode", wantErr: true, expectedCoupon: nil},
	}

	mockRepo := &MockRepository{
		coupons: map[string]entity.Coupon{
			"BlackFriday2024": {
				ID:             "00001",
				Code:           "BlackFriday2024",
				Discount:       10,
				MinBasketValue: 30,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Service{repo: mockRepo}
			gotCoupon, err := s.GetCoupon(tt.code)
			t.Logf("gotCoupon: -%v-", gotCoupon)

			if tt.wantErr {
				t.Logf("gotCoupon: -%v-", gotCoupon)
				assert.Error(t, err, "Expected an error for coupon code %s, but got none", tt.code)
				assert.Nil(t, gotCoupon, "Expected no coupon to be returned for invalid code %s, but got one", tt.code)
			} else {
				assert.NoError(t, err, "Expected no error for coupon code %s, but got one", tt.code)
				assert.NotNil(t, gotCoupon, "Expected a coupon for valid code %s, but got nil", tt.code)
				assert.Equal(t, tt.expectedCoupon, gotCoupon, "Coupon details do not match for code %s", tt.code)
			}
		})

	}
}
