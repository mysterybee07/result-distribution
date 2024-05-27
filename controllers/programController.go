package controllers

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/mysterybee07/result-distribution-system/initializers"
	"github.com/mysterybee07/result-distribution-system/models"
)

func AddProgram(c *fiber.Ctx) error {
	err := c.Render("dashboard/program/program", fiber.Map{})
	if err != nil {
		c.Status(fiber.StatusInternalServerError).SendString("Error rendering page")
		return err
	}
	return nil
}

func StoreProgram(c *fiber.Ctx) error {
	program := new(models.Program)

	// Parse the form data
	if err := c.BodyParser(program); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse form data",
		})
	}

	// Log the parsed data for debugging
	log.Printf("Parsed Program: %+v\n", program)

	var existingProgram models.Program
	// Check if the program already exists
	if err := initializers.DB.Where("name = ?", program.Name).First(&existingProgram).Error; err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Program already exists",
		})
	}

	// Create the new program
	if err := initializers.DB.Create(&program).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not create program",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Program created successfully",
		"program": program,
	})
}
