package middleware

import (
	"image-service/internal/config"

	"github.com/gofiber/fiber/v2"
)

func Auth(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		apiKey := c.Get("x-api-key")
		if apiKey == "" || apiKey != cfg.SecretKey {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized: Invalid API key",
			})
		}
		return c.Next()
	}
}
