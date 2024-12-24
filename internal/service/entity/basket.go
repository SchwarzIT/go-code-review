package entity

import (
	_ "github.com/gin-gonic/gin"
)

type Basket struct {
	Value                 int
	AppliedDiscount       int
	ApplicationSuccessful bool
}
