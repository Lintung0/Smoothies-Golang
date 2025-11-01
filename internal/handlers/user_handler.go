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

// CreateOrder handles the creation of a new order.
func (h *UserHandler) CreateOrder(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	var input models.Order
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	// Set UserID from JWT
	input.UserID = userID
	input.Status = "pending" // Default status

	// Start a database transaction
	tx := h.DB.Begin()
	if tx.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to start transaction"})
	}

	// Create the order
	// We create the order shell first, without items, to get an ID.
	// The total will be updated later.
	order := models.Order{
		UserID:           input.UserID,
		Status:           input.Status,
		NamaPenerima:     input.NamaPenerima,
		KelasAlamat:      input.KelasAlamat,
		Catatan:          input.Catatan,
		MetodePembayaran: input.MetodePembayaran,
	}
	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create order"})
	}

	var total float64 = 0
	for i := range input.OrderItems {
		// 1. Find product and check stock
		var product models.Products
		if err := tx.First(&product, input.OrderItems[i].ProductID).Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Product not found"})
		}
		if product.Stock < input.OrderItems[i].Qty {
			tx.Rollback()
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Not enough stock for product " + product.Name})
		}

		// 2. Update stock
		product.Stock -= input.OrderItems[i].Qty
		if err := tx.Save(&product).Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update stock"})
		}
	}

	// 3. Calculate total and prepare order items for creation
	for i := range input.OrderItems {
		var product models.Products
		// This read is safe as it's within the transaction
		tx.First(&product, input.OrderItems[i].ProductID)
		input.OrderItems[i].Price = product.Price
		input.OrderItems[i].OrderID = order.ID
		total += product.Price * float64(input.OrderItems[i].Qty)
	}

	// 4. Create all order items at once
	if err := tx.Create(&input.OrderItems).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create order items"})
	}

	// 5. Update the order with the final total
	order.Total = total
	if err := tx.Save(&order).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update order total"})
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to commit transaction"})
	}

	// Return the complete order object
	order.OrderItems = input.OrderItems
	return c.Status(fiber.StatusCreated).JSON(order)
}

// GetUserOrders retrieves the order history for a user.
func (h *UserHandler) GetUserOrders(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	var orders []models.Order
	h.DB.Preload("OrderItems.Product").Where("user_id = ?", userID).Find(&orders)
	return c.JSON(orders)
}

// Get profile by ID
func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint) // diambil dari middleware JWT

	var user models.User
	if err := h.DB.First(&user, userID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "user not found"})
	}

	return c.JSON(user)
}

// Update profile
func (h *UserHandler) UpdateProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	var user models.User
	if err := h.DB.First(&user, userID).Error; err != nil {
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
