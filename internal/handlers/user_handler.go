package handlers

import (
	"github.com/Lintung0/Smoothies-Golang/internal/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type UserHandler struct {
	DB *gorm.DB
}

func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{DB: db}
}

// GetUserOrders retrieves the order history for a user.
func (h *UserHandler) GetUserOrders(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	var orders []models.Order
	h.DB.Preload("OrderItems.Product").Where("user_id = ?", userID).Order("created_at desc").Find(&orders)
	return c.JSON(orders)
}

// Get profile by ID
func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint) // diambil dari middleware JWT

	var user models.User
	if err := h.DB.First(&user, "id = ?", userID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "user not found"})
	}

	return c.JSON(user)
}

// Update profile
func (h *UserHandler) UpdateProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	var user models.User
	if err := h.DB.First(&user, "id = ?", userID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "user not found"})
	}

	var input struct {
		Name   string `json:"name"`
		Kelas  string `json:"kelas"`
		Alamat string `json:"alamat"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	user.Name = input.Name
	user.Kelas = input.Kelas
	user.Alamat = input.Alamat

	if err := h.DB.Save(&user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(user)
}
