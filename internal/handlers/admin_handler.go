package handlers

import (
	"github.com/Lintung0/Smoothies-Golang/internal/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type AdminHandler struct {
	DB *gorm.DB
}

func NewAdminHandler(db *gorm.DB) *AdminHandler {
	return &AdminHandler{DB: db}
}

// GetOrders handles fetching all orders with pagination
func (h *AdminHandler) GetOrders(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("pageSize", 10)

	var orders []models.Order
	offset := (page - 1) * pageSize

	if err := h.DB.Preload("User").Preload("OrderItems.Product").Preload("Payment").
		Limit(pageSize).Offset(offset).Find(&orders).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"page":     page,
		"pageSize": pageSize,
		"orders":   orders,
	})
}

// UpdateOrderStatus handles updating the status of an order
func (h *AdminHandler) UpdateOrderStatus(c *fiber.Ctx) error {
	id := c.Params("id")
	var order models.Order

	// Find the order by ID
	if err := h.DB.First(&order, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "order not found"})
	}

	// Parse the new status from the request body
	var body struct {
		Status string `json:"status"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request body"})
	}

	// Update the order status
	order.Status = body.Status
	if err := h.DB.Save(&order).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "order status updated", "order": order})
}

// GetWeeklySales handles fetching weekly sales statistics
func (h *AdminHandler) GetWeeklySales(c *fiber.Ctx) error {
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")

	if startDate == "" || endDate == "" {
		return c.Status(400).JSON(fiber.Map{"error": "startDate and endDate are required"})
	}

	var sales []struct {
		Date  string  `json:"date"`
		Total float64 `json:"total"`
	}

	query := `
		SELECT DATE(created_at) as date, SUM(total) as total
		FROM orders 
		WHERE created_at BETWEEN ? AND ?
		GROUP BY DATE(created_at)
		ORDER BY DATE(created_at)
	`

	if err := h.DB.Raw(query, startDate, endDate).Scan(&sales).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"sales": sales})
}

// Get all users
func (h *AdminHandler) GetAllUsers(c *fiber.Ctx) error {
	var users []models.User
	if err := h.DB.Find(&users).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(users)
}

// Verify payment (admin only)
func (h *AdminHandler) VerifyPayment(c *fiber.Ctx) error {
	id := c.Params("id")
	var payment models.Payment
	if err := h.DB.First(&payment, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "payment not found"})
	}

	payment.Verified = true
	if err := h.DB.Save(&payment).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "payment verified", "payment": payment})
}
