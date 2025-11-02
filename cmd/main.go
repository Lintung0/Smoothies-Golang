package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"github.com/Lintung0/Smoothies-Golang/internal/config"
	"github.com/Lintung0/Smoothies-Golang/internal/handlers"
	"github.com/Lintung0/Smoothies-Golang/internal/middleware"
)

func main() {
	// 🔧 Load environment & connect to DB
	config.LoadEnv()
	db := config.ConnectDB()

	// 🚀 Fiber instance
	app := fiber.New()
	app.Use(logger.New())
	app.Use(middleware.SetupCors())

	// 📂 Buat folder upload jika belum ada
	if _, err := os.Stat("./uploads"); os.IsNotExist(err) {
		_ = os.MkdirAll("./uploads/products", os.ModePerm)
		_ = os.MkdirAll("./uploads/payments", os.ModePerm)
	}
	app.Static("/uploads", "./uploads")

	// 🧩 Handler setup
	authHandler := handlers.NewAuthHandler(db)
	userHandler := handlers.NewUserHandler(db)
	adminHandler := handlers.NewAdminHandler(db)
	productHandler := handlers.NewProductHandler(db)
	orderHandler := handlers.NewOrderHandler(db)

	// 🌐 Base group
	api := app.Group("/api")

	// ================================
	// 🔓 Public Routes
	// ================================
	api.Post("/register", authHandler.Register)
	api.Post("/login", authHandler.Login)

	// Produk bisa dilihat semua orang
	api.Get("/products", productHandler.GetAll)
	api.Get("/products/:id", productHandler.GetByID)

	// ================================
	// 👤 User Routes (require user role)
	// ================================
	userRoutes := api.Group("/user", middleware.AuthRequired("user"))
	userRoutes.Get("/profile", userHandler.GetProfile)
	userRoutes.Put("/profile", userHandler.UpdateProfile)
	userRoutes.Post("/orders", orderHandler.Create)                         // Buat order dengan raw JSON
	userRoutes.Post("/orders/:id/payment", orderHandler.UploadPaymentProof) // Upload bukti bayar (multipart)
	userRoutes.Get("/orders", userHandler.GetUserOrders)                    // lihat riwayat pesanan

	// ================================
	// 🛒 Admin Routes (require admin role)
	// ================================
	adminRoutes := api.Group("/admin", middleware.AuthRequired("admin"))

	// CRUD produk
	adminRoutes.Post("/products", productHandler.Create)
	adminRoutes.Put("/products/:id", productHandler.Update)
	adminRoutes.Delete("/products/:id", productHandler.Delete)

	// Order management
	adminRoutes.Get("/orders", adminHandler.GetOrders)                    // semua pesanan (pagination)
	adminRoutes.Put("/orders/:id/status", adminHandler.UpdateOrderStatus) // ubah status pesanan

	// User management
	adminRoutes.Get("/users", adminHandler.GetAllUsers) // lihat semua users

	// Payment verification
	adminRoutes.Put("/payments/:id/verify", adminHandler.VerifyPayment) // verify payment

	// Statistik penjualan mingguan
	adminRoutes.Get("/sales/weekly", adminHandler.GetWeeklySales)

	// ================================
	// 🚀 Run server
	// ================================
	port := config.GetEnv("APP_PORT", "8080")
	log.Printf("✅ Server running on port :%s\n", port)
	log.Fatal(app.Listen(":" + port))
}
