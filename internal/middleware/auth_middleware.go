package middleware

import (
	"github.com/Lintung0/Smoothies-Golang/internal/utils"

	"github.com/gofiber/fiber/v2"
)

func AuthRequired(role string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Get("Authorization")
		if token == "" {
			return c.Status(401).JSON(fiber.Map{"error": "Missing token"})
		}

		user, err := utils.ParseJWT(token)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "Invalid token"})
		}

		if role != "" && user.Role != role {
			return c.Status(403).JSON(fiber.Map{"error": "Forbidden"})
		}

		c.Locals("user", user)
		c.Locals("user_id", user.ID) // Tambahkan ini
		return c.Next()
	}
}
