package entity

import "coupon_service/internal/service/entity"

type ApplicationRequest struct {
	Code   string        `json:"code" binding:"required"`
	Basket entity.Basket `json:"basket" binding:"required"`
}
