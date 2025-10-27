package db

import (
	"fmt"
	"log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"github.com/Lintung0/Smoothies-Golang/internal/models"
	"github.com/Lintung0/Smoothies-Golang/internal/configs"
)

var DB *gorm.DB

func connect() {
	dsn := configs.GetDBSn
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
	  log.Fatal("Failed to connect to database:", err)
	}
	DB=db

	DB.AutoMigrate(&models.User{}, &models.Product{}, &models.Order{}, &models.OrderItem{}, &models.Payment{})
	fmt.Println("DB connected and migrated")
}