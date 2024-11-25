package entity

type Basket struct {
	Value                 int  `json:"value,omitempty"`
	AppliedDiscount       int  `json:"applied_discount,omitempty"`
	ApplicationSuccessful bool `json:"application_successful,omitempty"`
}
