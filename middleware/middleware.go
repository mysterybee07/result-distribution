package middleware

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/mysterybee07/result-distribution-system/utils"
)

func AuthRequired(c *fiber.Ctx) error {
	// Get token from cookies
	token := c.Cookies("jwt")

	if token == "" {
		return c.Redirect("/login", fiber.StatusFound)
	}

	// Validate token
	userID, err := utils.ParseJwt(token)
	if err != nil {
		log.Printf("Failed to parse JWT: %v\n", err)
		// return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		return c.Redirect("/login", fiber.StatusFound)
	}

	log.Printf("Authenticated user ID: %s\n", userID)

	// Set user ID in locals
	c.Locals("userID", userID)
	return c.Next()
}
