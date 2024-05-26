package controllers

import "github.com/gofiber/fiber/v2"

func Index(c *fiber.Ctx) error {
	err := c.Render("dashboard/index", fiber.Map{})
	if err != nil {
		c.Status(fiber.StatusInternalServerError).SendString("Error rendering page")
		return err
	}
	return nil
}

// func AddStudent(c *fiber.Ctx) error {

// }
