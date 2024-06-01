package controllers

import "github.com/gofiber/fiber/v2"

func NotFound(c *fiber.Ctx) error {
	err := c.Render("error/error404", fiber.Map{})
	if err != nil {
		c.Status(fiber.StatusInternalServerError).SendString("Error rendering page")
		return err
	}
	return nil
}

func ServerError(c *fiber.Ctx) error {
	err := c.Render("error/error500", fiber.Map{})
	if err != nil {
		c.Status(fiber.StatusInternalServerError).SendString("Error rendering page")
		return err
	}
	return nil
}
