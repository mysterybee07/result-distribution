package middleware

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/mysterybee07/result-distribution-system/initializers"
	"github.com/mysterybee07/result-distribution-system/models"
	"github.com/mysterybee07/result-distribution-system/utils"
)

func AuthRequired(c *fiber.Ctx) error {
	// Get token from cookies
	token := c.Cookies("jwt")

	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "unauthorized",
		})
	}

	// Validate token and get user ID and role
	userID, _, err := utils.ParseJwt(token)
	if err != nil {
		log.Printf("Failed to parse JWT: %v\n", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "unauthorized user",
		})
	}

	// log.Printf("Authenticated user ID: %s, Role: %s\n", userID, role)

	// Set user ID and role in locals
	c.Locals("userID", userID)

	return c.Next()
}

func AdminRequired(c *fiber.Ctx) error {
	// Get user ID from locals
	userID := c.Locals("userID")
	if userID == nil {
		log.Println("User ID not found in locals")
		return c.Redirect("/login", fiber.StatusFound)
	}

	// Fetch user from database to check role
	var user models.User
	if err := initializers.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		log.Printf("Failed to fetch user: %v\n", err)
		return c.Redirect("/login", fiber.StatusFound)
	}

	if user.Role != "admin" && user.Role != "superadmin" {
		log.Println("User is not an admin")
		return c.Redirect("/login", fiber.StatusFound)
	}

	return c.Next()
}

func SuperadminRequired(c *fiber.Ctx) error {
	// Get user ID from locals
	userID := c.Locals("userID")
	if userID == nil {
		log.Println("User ID not found in locals")
		return c.Redirect("/login", fiber.StatusFound)
	}

	// Fetch user from the database to check role
	var user models.User
	if err := initializers.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		log.Printf("Failed to fetch user: %v\n", err)
		return c.Redirect("/login", fiber.StatusFound)
	}

	// Check if the user is a superadmin
	if user.Role != "superadmin" {
		log.Println("User is not a superadmin")
		// Redirect to another route (404 page or dashboard)
		return c.Redirect("/404", fiber.StatusFound)
	}

	return c.Next()
}

// var store = session.New()

// func FlashMessages(c *fiber.Ctx) error {
// 	sess, err := store.Get(c)
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).SendString("Failed to get session")
// 	}

// 	if msg := sess.Get("success"); msg != nil {
// 		c.Locals("flash_success", msg)
// 		sess.Delete("success")
// 		if err := sess.Save(); err != nil {
// 			return c.Status(fiber.StatusInternalServerError).SendString("Failed to save session")
// 		}
// 	}

// 	return c.Next()
// }
