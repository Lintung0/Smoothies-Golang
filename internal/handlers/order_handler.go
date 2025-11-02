package handlers

import (
	"fmt"
	"path/filepath"

	"github.com/Lintung0/Smoothies-Golang/internal/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type OrderHandler struct {
	DB *gorm.DB
}

func NewOrderHandler(db *gorm.DB) *OrderHandler {
	return &OrderHandler{DB: db}
}

// Create handles the creation of a new order.
func (h *OrderHandler) Create(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	// Menggunakan BodyParser untuk membaca raw JSON dari request body.
	var input models.Order
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	// Validate that there is at least one order item
	if len(input.OrderItems) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Order must have at least one item"})
	}

	// Start a database transaction
	tx := h.DB.Begin()
	if tx.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to start transaction"})
	}

	// Create the order
	order := models.Order{
		UserID:           userID,    // Gunakan userID dari token, bukan dari input
		Status:           "pending", // Default status
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
	// Buat slice baru yang bersih untuk menampung item yang akan dibuat.
	// Ini menghindari modifikasi langsung pada data input.
	itemsToCreate := make([]models.OrderItem, 0, len(input.OrderItems))
	// Process, validate, and prepare order items in a single loop
	for i := range input.OrderItems {
		var product models.Products
		// 1. Find product
		if err := tx.First(&product, input.OrderItems[i].ProductID).Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Product not found"})
		}

		// 2. Check stock
		if product.Stock < input.OrderItems[i].Qty {
			tx.Rollback()
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Not enough stock for product " + product.Name})
		}

		// 3. Update stock
		product.Stock -= input.OrderItems[i].Qty
		if err := tx.Save(&product).Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update stock"})
		}

		// 4. Prepare order item for creation
		item := models.OrderItem{
			OrderID:   order.ID,
			ProductID: input.OrderItems[i].ProductID,
			Qty:       input.OrderItems[i].Qty,
			Price:     product.Price,
		}
		itemsToCreate = append(itemsToCreate, item)
		total += product.Price * float64(input.OrderItems[i].Qty)
	}

	// 5. Create all order items at once
	if err := tx.Create(&itemsToCreate).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create order items"})
	}

	// 6. Update the order with the final total and associate the created items for the response
	order.Total = total
	order.OrderItems = itemsToCreate // Gunakan slice yang baru dibuat yang sekarang berisi ID dari database
	if err := tx.Save(&order).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update order total"})
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to commit transaction"})
	}

	return c.Status(fiber.StatusCreated).JSON(order)
}

// UploadPaymentProof handles uploading a payment proof for an existing order.
func (h *OrderHandler) UploadPaymentProof(c *fiber.Ctx) error {
	orderID := c.Params("id")

	// Handle file upload for payment proof
	file, err := c.FormFile("payment_proof")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing 'payment_proof' file in form"})
	}

	// Buat nama file yang unik untuk menghindari konflik
	filename := fmt.Sprintf("payment-%s-%s", orderID, filepath.Base(file.Filename))
	filePath := filepath.Join("./uploads/payments", filename)

	// Simpan file ke disk
	if err := c.SaveFile(file, filePath); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save payment proof"})
	}

	// Buat record payment di database
	payment := models.Payment{
		PaymentProofURL: "/" + filepath.ToSlash(filePath), // Simpan path yang bisa diakses web
	}
	// Konversi orderID dari string ke uint
	if _, err := fmt.Sscan(orderID, &payment.OrderID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid order ID"})
	}

	if err := h.DB.Create(&payment).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create payment record"})
	}

	return c.Status(fiber.StatusCreated).JSON(payment)
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
