package api

import (
	"bytes"
	. "coupon_service/internal/api/entity"
	"coupon_service/internal/service/entity"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockCouponService struct {
	mock.Mock
}

func (m *mockCouponService) ApplyCoupon(basket entity.Basket, code string) (*entity.Basket, error) {
	args := m.Called(basket, code)
	return args.Get(0).(*entity.Basket), args.Error(1)
}

func (m *mockCouponService) CreateCoupon(discount int, code string, minBasketValue int) error {
	args := m.Called(discount, code, minBasketValue)
	return args.Error(0)
}

func (m *mockCouponService) GetCoupons(codes []string) ([]entity.Coupon, error) {
	args := m.Called(codes)
	return args.Get(0).([]entity.Coupon), args.Error(1)
}

func setupRouter(api *API) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/apply", api.Apply)
	router.POST("/coupons", api.Create)
	router.POST("/get", api.Get)
	return router
}

func TestAPI_Apply(t *testing.T) {
	tests := []struct {
		name         string
		requestBody  ApplicationRequest
		expectedCode int
		expectedBody *entity.Basket
		mockReturn   *entity.Basket
		mockError    error
	}{
		{
			name: "apply valid coupon",
			requestBody: ApplicationRequest{
				Basket: entity.Basket{Value: 100},
				Code:   "VALID",
			},
			expectedCode: http.StatusOK,
			expectedBody: &entity.Basket{
				Value:                 100,
				AppliedDiscount:       10,
				ApplicationSuccessful: true,
			},
			mockReturn: &entity.Basket{
				Value:                 100,
				AppliedDiscount:       10,
				ApplicationSuccessful: true,
			},
			mockError: nil,
		},
		{
			name: "apply invalid coupon",
			requestBody: ApplicationRequest{
				Basket: entity.Basket{Value: 100},
				Code:   "INVALID",
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: nil,
			mockReturn:   nil,
			mockError:    fmt.Errorf("invalid coupon code"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			couponService := &mockCouponService{}
			couponService.On("ApplyCoupon", tt.requestBody.Basket, tt.requestBody.Code).Return(tt.mockReturn, tt.mockError)
			api := New(Config{}, couponService)
			router := setupRouter(&api)

			body, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest(http.MethodPost, "/apply", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
			if tt.expectedBody != nil {
				var responseBody entity.Basket
				err := json.Unmarshal(w.Body.Bytes(), &responseBody)
				assert.NoError(t, err)
				assert.Equal(t, *tt.expectedBody, responseBody)
			}
		})
	}
}

func TestAPI_Create(t *testing.T) {
	tests := []struct {
		name         string
		requestBody  entity.Coupon
		expectedCode int
		mockError    error
	}{
		{
			name: "create valid coupon",
			requestBody: entity.Coupon{
				Discount:       10,
				Code:           "NEWCOUPON",
				MinBasketValue: 50,
			},
			expectedCode: http.StatusOK,
			mockError:    nil,
		},
		{
			name: "create duplicate coupon",
			requestBody: entity.Coupon{
				Discount:       10,
				Code:           "DUPLICATE",
				MinBasketValue: 50,
			},
			expectedCode: http.StatusBadRequest,
			mockError:    fmt.Errorf("coupon code already exists"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			couponService := &mockCouponService{}
			couponService.On("CreateCoupon", tt.requestBody.Discount, tt.requestBody.Code, tt.requestBody.MinBasketValue).Return(tt.mockError)
			api := New(Config{}, couponService)
			router := setupRouter(&api)

			body, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest(http.MethodPost, "/coupons", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
		})
	}
}

func TestAPI_Get(t *testing.T) {
	tests := []struct {
		name         string
		requestBody  CouponRequest
		expectedCode int
		expectedBody []entity.Coupon
		mockReturn   []entity.Coupon
		mockError    error
	}{
		{
			name: "get existing coupons",
			requestBody: CouponRequest{
				Codes: []string{"VALID"},
			},
			expectedCode: http.StatusOK,
			expectedBody: []entity.Coupon{
				{
					Code:           "VALID",
					Discount:       10,
					MinBasketValue: 50,
				},
			},
			mockReturn: []entity.Coupon{
				{
					Code:           "VALID",
					Discount:       10,
					MinBasketValue: 50,
				},
			},
			mockError: nil,
		},
		{
			name: "get non-existing coupon",
			requestBody: CouponRequest{
				Codes: []string{"INVALID"},
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: nil,
			mockReturn:   nil,
			mockError:    fmt.Errorf("invalid coupon code"),
		},
		{
			name: "get mixed existing and non-existing coupons",
			requestBody: CouponRequest{
				Codes: []string{"VALID", "INVALID"},
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: nil,
			mockReturn:   nil,
			mockError:    fmt.Errorf("invalid coupon code"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			couponService := &mockCouponService{}
			couponService.On("GetCoupons", tt.requestBody.Codes).Return(tt.mockReturn, tt.mockError)
			api := New(Config{}, couponService)
			router := setupRouter(&api)

			body, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest(http.MethodPost, "/get", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
			if tt.expectedBody != nil {
				var responseBody []entity.Coupon
				err := json.Unmarshal(w.Body.Bytes(), &responseBody)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBody, responseBody)
			}
		})
	}
}
