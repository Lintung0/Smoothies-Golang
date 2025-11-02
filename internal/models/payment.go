package models

import "gorm.io/gorm"

type Payment struct {
	gorm.Model
	OrderID         uint    `json:"order_id"`
	ImageURL        string  `json:"image_url"`
	PaymentProofURL string  `json:"payment_proof_url"`
	Verified        bool    `json:"verified" gorm:"default:false"`
	Amount          float64 `json:"amount"`
}
