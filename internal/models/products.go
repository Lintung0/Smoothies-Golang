package models

import "gorm.io/gorm"

type Products struct {
	gorm.Model
	Name        string      `json:"name"`
	Description string      `json:"description" gorm:"type:text"`
	Price       float64     `json:"price"`
	Stock       int         `json:"stock"`
	ImageURL    string      `json:"image_url"`
	OrderItems  []OrderItem `gorm:"foreignKey:ProductID"`
}
