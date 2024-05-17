package controllers

import "github.com/gofiber/fiber/v2"

func Profile(c *fiber.Ctx) error {
	err := c.Render("users/profile", fiber.Map{})
	if err != nil {
		c.Status(fiber.StatusInternalServerError).SendString("Error rendering page")
		return err
	}
	return nil
}
