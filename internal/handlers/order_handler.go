package handlers

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"github.com/Lintung0/Smoothies-Golang/internal/models"
)

type OrderHandler struct {
	DB *gorm.DB
}

func NewOrderHandler(db *gorm.DB) *OrderHandler {
	return &OrderHandler{DB: db}
}

// Create new order
func (h *OrderHandler) Create(c *fiber.Ctx) error {
	var order models.Order
	if err := c.BodyParser(&order); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	if err := h.DB.Create(&order).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(201).JSON(order)
}

// Get all orders
func (h *OrderHandler) GetAll(c *fiber.Ctx) error {
	var orders []models.Order
	if err := h.DB.Preload("OrderItems.Product").Preload("Payment").Find(&orders).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(orders)
}

// Get order by ID
func (h *OrderHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var order models.Order
	if err := h.DB.Preload("OrderItems.Product").Preload("Payment").First(&order, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "order not found"})
	}
	return c.JSON(order)
}

// Update order status
func (h *OrderHandler) UpdateStatus(c *fiber.Ctx) error {
	id := c.Params("id")
	var order models.Order
	if err := h.DB.First(&order, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "order not found"})
	}

	type StatusUpdate struct {
		Status string `json:"status"`
	}
	var input StatusUpdate
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	order.Status = input.Status
	if err := h.DB.Save(&order).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(order)
}
