package config

import (
	"fmt"
	"log"

	"github.com/Lintung0/Smoothies-Golang/internal/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectDB() *gorm.DB {
	LoadEnv()

	user := GetEnv("DB_USER", "Lintang")
	pass := GetEnv("DB_PASS", "Lintang160309")
	host := GetEnv("GB_HOST", "127.0.0.1")
	port := GetEnv("DB_PORT", "3306")
	dbname := GetEnv("DB_NAME", "smoothiesdb")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, pass, host, port, dbname)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ Failed to connect to database:", err)
	}

	log.Println("✅ Connected to database")
	db.AutoMigrate(&models.User{}, &models.Products{}, &models.Order{}, &models.OrderItem{}, &models.Payment{})
	return db
}
