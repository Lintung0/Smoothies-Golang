package handlers

import (
	"encoding/json"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/Lintung0/Smoothies-Golang/internal/db"
	"github.com/Lintung0/Smoothies-Golang/internal/models"
)

func CreateOrder(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	// ambil data dari form
	nama := c.FormValue("nama_penerima")
	kelasAlamat := c.FormValue("kelas_alamat")
	catatan := c.FormValue("catatan")
	metode := c.FormValue("metode_pembayaran")
	itemsJSON := c.FormValue("items")

	var items []models.OrderItem
	if err := json.Unmarshal([]byte(itemsJSON), &items); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Format item salah"})
	}

	var total float64
	for _, i := range items {
		total += i.Price * float64(i.Qty)
	}

	order := models.Order{
		UserID:           uint64(userID),
		Total:            total,
		NamaPenerima:     nama,
		KelasAlamat:      kelasAlamat,
		Catatan:          catatan,
		MetodePembayaran: metode,
		Status:           "pending",
		CreatedAt:        time.Now(),
	}

	db.DB.Create(&order)

	for _, i := range items {
		i.OrderID = order.ID
		db.DB.Create(&i)
	}

	// upload bukti pembayaran (optional)
	file, err := c.FormFile("payment_proof")
	if err == nil {
		path := "uploads/" + file.Filename
		c.SaveFile(file, path)
		payment := models.Payment{
			OrderID:         order.ID,
			PaymentProofURL: path,
			Verified:        false,
			Amount:          total,
		}
		db.DB.Create(&payment)
	}

	return c.JSON(fiber.Map{"message": "Pesanan dibuat", "order_id": order.ID})
}

func GetUserOrders(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	var orders []models.Order
	db.DB.Preload("OrderItems").Preload("Payments").
		Where("user_id = ?", userID).Order("created_at desc").
		Find(&orders)
	return c.JSON(orders)
}
