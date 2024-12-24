package entity

import "coupon_service/internal/service/entity"

type ApplicationRequest struct {
	Code   string
	Basket entity.Basket
}
