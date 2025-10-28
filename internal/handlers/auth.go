package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/Lintung0/Smoothies-Golang/internal/db"
	"github.com/Lintung0/Smoothies-Golang/internal/models"
	"os"
)

func Register(c *fiber.Ctx) error {
	var req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Kelas    string `json:"kelas"`
		Alamat   string `json:"alamat"`
	}

	if err := c.BodyParser(&req); err != nil {
	  return c.Status(400).JSON(fiber.Map{"error": "Data tidak valid"})
	}

	hash, _:= bcrypt.GenerateFromPassword([]byte(req.Password), 14)

	user := models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hash),
		Role:     "user",
		Kelas:    req.Kelas,
		Alamat:   req.Alamat,
	}

	if err := db.DB.Create(&user).Error; err != nil {
	  return c.Status(500).JSON(fiber.Map{"error": "Gagal mendaftarkan pengguna"})
	}

	return c.JSON(fiber.Map{"message": "Registrasi berhasil"})
}

func Login(c *fiber.Ctx) error {
	var req struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&req); err !=nil {
	  return c.Status(400).JSON(fiber.Map{"error": "Format salah"})
	}

	var user models.User
	if err := db.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
	  return c.Status(401).JSON(fiber.Map{"error": "User tidak ditemukan"})
	}
	
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Password salah"})
	}

	claims := jwt.MapClaims{
		"sub":  user.ID,
		"role": user.Role,
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, _ := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	return c.JSON(fiber.Map{
		"access_token": signedToken,
		"role": user.Role,
	})
}
