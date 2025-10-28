package models

import (
	"time"
)

type User struct {
	ID		uint     `gorm:"primaryKey" json:"id"`
	Name 	string   `json:"name"`
	Email 	string   `gorm:"unique" json:"email"`
	Password string  `json:"-"`
	Role string   `json:"role"` // e.g., "admin" or "customer"
	Kelas string `json:"kelas"`
	Alamat string `json:"alamat"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Products struct {
	ID        uint64 `gorm:"primaryKey" json:"id"`
	Name      string `json:"name"`
	Description string `json:"description"`
	Price     float64 `json:"price"`
	Stock     int     `json:"stock"`
	ImageURL  string  `json:"image_url"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Order struct {
	ID uint64 `gorm:"primaryKey" json:"id"`
	UserID uint64 `json:"user_id"`
	Total float64 `json:"total"`
	NamaPenerima string `json:"nama_penerima"`
	KelasAlamat string `json:"kelas_alamat"`
	Catatan string `json:"catatan"`
	MetodePembayaran string `json:"metode_pembayaran"`
	Status string `json:"status"`
	OrderItems []OrderItem `gorm:"foreignKey:OrderID" json:"items"`
	Payments []Payment `gorm:"foreignKey:OrderID" json:"payments"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type OrderItem struct {
	ID uint64 `gorm:"primaryKey" json:"id"`
	OrderID uint64 `json:"order_id"`
	ProductID uint64 `json:"product_id"`
	Qty int `json:"qty"`
	Price float64 `json:"price"`
	CreatedAt time.Time
}

type Payment struct {
	ID uint64 `gorm:"primaryKey" json:"id"`
	OrderID uint64 `json:"order_id"`
	PaymentProofURL string `json:"payment_proof_url"`
	Verified bool `json:"verified"`
	Amount float64 `json:"amount"`
	CreatedAt time.Time
}