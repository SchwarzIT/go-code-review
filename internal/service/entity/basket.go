package entity

import (
	_ "github.com/gin-gonic/gin"
)

type Basket struct {
	Value                 int  `json:"Value" binding:"required"`
	AppliedDiscount       int  `json:"applied_discount,omitempty"`
	ApplicationSuccessful bool `json:"application_successful,omitempty"`
}
