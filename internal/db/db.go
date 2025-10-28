package db

import (
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/Lintung0/Smoothies-Golang/internal/models"
)

var DB *gorm.DB

func Connect() {
	dsn := os.Getenv("DB_DSN")
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ gagal konek database:", err)
	}

	err = DB.AutoMigrate(&models.User{}, &models.Products{}, &models.Order{}, &models.OrderItem{}, &models.Payment{})
	if err != nil {
		log.Fatal("❌ gagal migrasi tabel:", err)
	}

	log.Println("✅ Database terkoneksi & migrasi berhasil")
}
