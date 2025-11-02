package models

import "gorm.io/gorm"

type OrderItem struct {
	gorm.Model
	OrderID   uint     `json:"order_id"`
	ProductID uint     `json:"product_id"`
	Qty       int      `json:"qty"`
	Price     float64  `json:"price"`
	Product   Products `json:"product" gorm:"foreignKey:ProductID"`
}
