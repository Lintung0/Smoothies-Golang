package handlers

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"github.com/Lintung0/Smoothies-Golang/internal/models"
)

type ProductHandler struct {
	DB *gorm.DB
}

func NewProductHandler(db *gorm.DB) *ProductHandler {
	return &ProductHandler{DB: db}
}

// Get all products
func (h *ProductHandler) GetAll(c *fiber.Ctx) error {
	var products []models.Products
	if err := h.DB.Find(&products).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(products)
}

// Get product by ID
func (h *ProductHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var product models.Products
	if err := h.DB.First(&product, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "product not found"})
	}
	return c.JSON(product)
}

// Create product
func (h *ProductHandler) Create(c *fiber.Ctx) error {
	var product models.Products
	if err := c.BodyParser(&product); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	if err := h.DB.Create(&product).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(201).JSON(product)
}

// Update product
func (h *ProductHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	var product models.Products
	if err := h.DB.First(&product, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "product not found"})
	}

	if err := c.BodyParser(&product); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	if err := h.DB.Save(&product).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(product)
}

// Delete product
func (h *ProductHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.DB.Delete(&models.Products{}, id).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "product deleted"})
}
