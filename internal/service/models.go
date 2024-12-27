package service

type Basket struct {
	Value                 int  `json:"Value" binding:"required,gt=0"`
	AppliedDiscount       int  `json:"applied_discount,omitempty" swaggerignore:"true"`
	ApplicationSuccessful bool `json:"application_successful,omitempty" swaggerignore:"true"`
}
