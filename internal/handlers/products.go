package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/Lintung0/Smoothies-Golang/internal/db"
	"github.com/Lintung0/Smoothies-Golang/internal/models"
	"strconv"
)

//get all product

func ListProducts(c *fiber.Ctx) error {
	var products []models.Products
	db.DB.Find(&products)
	return c.JSON(products)
}

// get product by id
func GetProductByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var product models.Products
	if err := db.DB.First(&product, id).Error; err != nil {
	  return c.Status(404).JSON(fiber.Map{"error": "Produk tidak ditemukan"})
	}
	return c.JSON(product)
}

//create product (admin only)
func CreateProduct(c *fiber.Ctx) error {
	name := c.FormValue("name")
	description := c.FormValue("description")
	price, _ := strconv.ParseFloat(c.FormValue("price"), 64)
	stock, _ := strconv.Atoi(c.FormValue("stock"))

	file, err :=c.FormFile("image")
	var filePath string
	if err == nil {
		filePath = "uploads/" + file.Filename
		c.SaveFile(file, filePath)
	}

	product := models.Products{
		Name:        name,
		Description: description,
		Price:       price,
		Stock:       stock,
		ImageURL:    filePath,
	}

	db.DB.Create(&product)
	return c.JSON(product)
}

//update product
func UpdateProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	var product models.Products
	if err := db.DB.First(&product, id).Error; err != nil {
	  return c.Status(404).JSON(fiber.Map{"error": "Produk tidak ditemukan"})
	}
	
	name := c.FormValue("name")
	description := c.FormValue("description")
	price, _ := strconv.ParseFloat(c.FormValue("price"), 64)
	stock, _ := strconv.Atoi(c.FormValue("stock"))

	file, err := c.FormFile("image")
	if err == nil {
		filePath := "uploads/" + file.Filename
		c.SaveFile(file, filePath)
		product.ImageURL = filePath
	}

	product.Name = name
	product.Description = description
	product.Price = price
	product.Stock = stock

	db.DB.Save(&product)
	return c.JSON(product)
}

// delete product 
func DeleteProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	db.DB.Delete(&models.Products{}, id)
	return c.JSON(fiber.Map{"Message": "Produk dihapus"})
}