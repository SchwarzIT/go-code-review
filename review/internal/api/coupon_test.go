package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"coupon_service/internal/entity"
)

type MockService struct{}

func (m *MockService) ApplyCoupon(basket *entity.Basket, code string) error {
	if code == "valid" {
		basket.Value -= 10
		return nil
	}
	return errors.New("invalid coupon code")
}

func (m *MockService) CreateCoupon(discount int, code string, minBasketValue int) error {
	if code == "duplicate" {
		return errors.New("duplicated coupon code")
	}
	return nil
}

func (m *MockService) GetCoupons(codes []string) []entity.Coupon {
	if len(codes) == 0 {
		return nil
	}
	if len(codes) == 1 {
		if codes[0] == "invalid" {
			return []entity.Coupon{}
		}
	}
	return []entity.Coupon{{Code: "valid", Discount: 10, MinBasketValue: 30}}
}

func (m *MockService) GetCoupon(code string) (*entity.Coupon, error) {
	if code == "valid" {
		return &entity.Coupon{Code: "valid", Discount: 10, MinBasketValue: 30}, nil
	}
	return nil, errors.New("coupon not found")
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	return router
}

func TestApply(t *testing.T) {
	router := setupRouter()
	api := API{svc: &MockService{}}
	router.POST("/apply", api.Apply)

	tests := []struct {
		name       string
		request    entity.ApplicationRequest
		statusCode int
	}{
		{"Valid coupon", entity.ApplicationRequest{Basket: entity.Basket{Value: 100}, Code: "valid"}, http.StatusOK},
		{"Invalid coupon", entity.ApplicationRequest{Basket: entity.Basket{Value: 100}, Code: "invalid"}, http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.request)
			req, _ := http.NewRequest(http.MethodPost, "/apply", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)

			assert.Equal(t, tt.statusCode, resp.Code)
		})
	}
}

func TestCreate(t *testing.T) {
	router := setupRouter()
	api := API{svc: &MockService{}}
	router.POST("/create", api.Create)

	tests := []struct {
		name       string
		request    entity.Coupon
		statusCode int
	}{
		{"Valid creation", entity.Coupon{Code: "new", Discount: 10, MinBasketValue: 20}, http.StatusOK},
		{"Duplicated creation", entity.Coupon{Code: "duplicate", Discount: 10, MinBasketValue: 20}, http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.request)
			req, _ := http.NewRequest(http.MethodPost, "/create", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)

			assert.Equal(t, tt.statusCode, resp.Code)
		})
	}
}

func TestGet(t *testing.T) {
	router := setupRouter()
	api := API{svc: &MockService{}}
	router.POST("/coupons", api.Get)

	tests := []struct {
		name       string
		request    entity.CouponRequest
		statusCode int
		expected   map[string]interface{}
	}{
		{
			"Single valid code",
			entity.CouponRequest{Codes: []string{"valid"}},
			http.StatusOK,
			map[string]interface{}{
				"data": map[string]interface{}{
					"coupons": []interface{}{
						map[string]interface{}{
							"code":             "valid",
							"discount":         10.0,
							"min_basket_value": 30.0,
						},
					},
				},
			},
		},
		{
			"Invalid code",
			entity.CouponRequest{Codes: []string{"invalid"}},
			http.StatusOK,
			map[string]interface{}{
				"data": map[string]interface{}{
					"coupons": []interface{}{},
				},
			},
		},
		{
			"Multiple codes, some valid and some invalid",
			entity.CouponRequest{Codes: []string{"valid", "invalid"}},
			http.StatusOK,
			map[string]interface{}{
				"data": map[string]interface{}{
					"coupons": []interface{}{
						map[string]interface{}{
							"code":             "valid",
							"discount":         10.0,
							"min_basket_value": 30.0,
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.request)
			req, _ := http.NewRequest(http.MethodPost, "/coupons", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)

			assert.Equal(t, tt.statusCode, resp.Code)

			var responseBody map[string]interface{}
			err := json.Unmarshal(resp.Body.Bytes(), &responseBody)
			assert.NoError(t, err)

			t.Logf("GOT: %v", responseBody)
			assert.Equal(t, tt.expected, responseBody)
		})
	}
}
