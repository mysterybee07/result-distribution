package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mysterybee07/result-distribution-system/utils"
)

func AuthRequired(c *fiber.Ctx) error {
	// Get token from cookies
	token := c.Cookies("jwt")

	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	// Validate token
	userID, err := utils.ParseJwt(token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	// Set user ID in locals
	c.Locals("userID", userID)
	return c.Next()
}
