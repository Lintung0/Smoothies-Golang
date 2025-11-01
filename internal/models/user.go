package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string  `json:"name"`
	Email    string  `json:"email" gorm:"unique"`
	Password string  `json:"-"`
	Role     string  `json:"role" gorm:"type:enum('user','admin');default:'user'"`
	Kelas    string  `json:"kelas"`
	Alamat   string  `json:"alamat" gorm:"type:text"`
	Orders   []Order `gorm:"foreignKey:UserID"`
}
