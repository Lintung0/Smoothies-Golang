package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"github.com/Lintung0/Smoothies-Golang/internal/db"
	"github.com/Lintung0/Smoothies-Golang/internal/handlers"
	"github.com/Lintung0/Smoothies-Golang/internal/middleware"
)

func main() {
	// konek database
	db.Connect()

	app := fiber.New()
	app.Use(logger.New())

	// Serve folder upload bukti pembayaran & produk
	app.Static("/uploads", "./uploads")

	api := app.Group("/api")

	// --- AUTH ---
	api.Post("/register", handlers.Register)
	api.Post("/login", handlers.Login)

	// --- PRODUCT (Public) ---
	api.Get("/products", handlers.ListProducts)
	api.Get("/products/:id", handlers.GetProductByID)

	// --- USER (Protected) ---
	user := api.Group("/user", middleware.JWTProtected)
	user.Post("/orders", handlers.CreateOrder)
	user.Get("/orders", handlers.GetUserOrders)

	// --- ADMIN (Protected) ---
	admin := api.Group("/admin", middleware.JWTProtected, middleware.RequireAdmin)
	admin.Post("/products", handlers.CreateProduct)
	admin.Put("/products/:id", handlers.UpdateProduct)
	admin.Delete("/products/:id", handlers.DeleteProduct)
	admin.Get("/orders", handlers.AdminGetOrders)
	admin.Put("/orders/:id/status", handlers.AdminUpdateOrderStatus)
	admin.Get("/sales/weekly", handlers.GetWeeklySales)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(app.Listen(":" + port))
}
