package entity

type ApplicationRequest struct {
	Code   string `json:"code,omitempty"`
	Basket Basket `json:"basket,omitempty"`
}
