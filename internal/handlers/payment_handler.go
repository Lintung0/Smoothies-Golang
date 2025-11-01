package handlers

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"github.com/Lintung0/Smoothies-Golang/internal/models"
)

type PaymentHandler struct {
	DB *gorm.DB
}

func NewPaymentHandler(db *gorm.DB) *PaymentHandler {
	return &PaymentHandler{DB: db}
}

// Create payment (upload bukti)
func (h *PaymentHandler) Create(c *fiber.Ctx) error {
	var payment models.Payment
	if err := c.BodyParser(&payment); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	if err := h.DB.Create(&payment).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(201).JSON(payment)
}

// Get all payments
func (h *PaymentHandler) GetAll(c *fiber.Ctx) error {
	var payments []models.Payment
	if err := h.DB.Find(&payments).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(payments)
}

// Verify payment
func (h *PaymentHandler) Verify(c *fiber.Ctx) error {
	id := c.Params("id")
	var payment models.Payment
	if err := h.DB.First(&payment, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "payment not found"})
	}

	payment.Verified = true
	if err := h.DB.Save(&payment).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(payment)
}
