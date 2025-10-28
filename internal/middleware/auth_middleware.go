package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"os"
)

func JWTProtected(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(401).JSON(fiber.Map{"error": "Token tidak ditemukan"})
	}

	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil || !token.Valid {
		return c.Status(401).JSON(fiber.Map{"error": "Token tidak valid"})
	}

	claims := token.Claims.(jwt.MapClaims)
	c.Locals("user_id", uint(claims["sub"].(float64)))
	c.Locals("role", claims["role"].(string))
	return c.Next()
}

func RequireAdmin(c *fiber.Ctx) error {
	role := c.Locals("role")
	if role != "admin" {
		return c.Status(403).JSON(fiber.Map{"error": "Akses khusus admin"})
	}
	return c.Next()
}
