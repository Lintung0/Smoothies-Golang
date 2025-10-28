package handlers

import (
	"strconv"
	"github.com/gofiber/fiber/v2"
	"github.com/Lintung0/Smoothies-Golang/internal/db"
	"github.com/Lintung0/Smoothies-Golang/internal/models"
)

func AdminGetOrders(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset := (page - 1) * limit

	var total int64
	db.DB.Model(&models.Order{}).Count(&total)

	var orders []models.Order
	db.DB.Preload("OrderItems").Preload("Payments").
	Order("created_at desc").Limit(limit).Offset(offset).Find(&orders)

	return c.JSON(fiber.Map{
		"data": orders,
		"page": page,
		"limit": limit,
		"total": total,
	})
}

func AdminUpdateOrderStatus(c *fiber.Ctx) error {
	id := c.Params("id")
	status := c.FormValue("status")
	db.DB.Model(&models.Order{}).Where("id = ?", id).Update("status", status)
	return c.JSON(fiber.Map{"Message":"Status pesanan diperbarui"})
}

// grafik penjualan tiap mingguan
func GetWeeklySales(c *fiber.Ctx) error {
	type Result struct {
		Year        int
		Week        int
		TotalSales  float64
	} 
	var results []Result
	db.DB.Raw(`
		SELECT YEAR(created_at) as year, WEEK(created_at,1) as week, SUM(total) as total_sales
		FROM orders
		GROUP BY year, week
		ORDER BY year, week
	`).Scan(&results)
	return c.JSON(results)
}