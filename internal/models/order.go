package models

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	User             User        `gorm:"foreignKey:UserID"`
	UserID           uint        `json:"user_id"`
	Total            float64     `json:"total"`
	NamaPenerima     string      `json:"nama_penerima"`
	KelasAlamat      string      `json:"kelas_alamat"`
	Catatan          string      `json:"catatan"`
	MetodePembayaran string      `json:"metode_pembayaran" gorm:"type:enum('qris','transfer','cod','lainnya')"`
	Status           string      `json:"status" gorm:"type:enum('pending','processing','completed','cancelled');default:'pending'"`
	OrderItems       []OrderItem `gorm:"foreignKey:OrderID"`
	Payment          Payment     `gorm:"foreignKey:OrderID"`
}
