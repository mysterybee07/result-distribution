package controllers

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/mysterybee07/result-distribution-system/initializers"
	"github.com/mysterybee07/result-distribution-system/models"
)

func Program(c *fiber.Ctx) error {
	// Fetch all programs from the database
	var programs []models.Program
	if err := initializers.DB.Find(&programs).Error; err != nil {
		log.Printf("Failed to fetch programs: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to fetch programs")
	}

	// Render the form with the list of programs
	err := c.Render("dashboard/program/program", fiber.Map{
		"Programs": programs,
	})
	if err != nil {
		log.Printf("Failed to render page: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Error rendering page")
	}

	return nil
}

func CreateProgram(c *fiber.Ctx) error {
	var program models.Program

	// Parse the form data
	if err := c.BodyParser(&program); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse form data",
		})
	}

	// Log the parsed data for debugging
	log.Printf("Parsed Program: %+v\n", program)

	var existingProgram models.Program
	// Check if the program already exists
	if err := initializers.DB.Where("program_name = ?", program.ProgramName).First(&existingProgram).Error; err == nil {
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

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Programs created successfully",
		"program": program,
	})
}

func EditProgram(c *fiber.Ctx) error {
	err := c.Render("dashboard/program/editprogram", fiber.Map{})
	if err != nil {
		c.Status(fiber.StatusInternalServerError).SendString("Error rendering page")
		return err
	}
	return nil
}

func UpdateProgram(c *fiber.Ctx) error {
	id := c.Params("id")

	var program models.Program

	if err := initializers.DB.First(&program, id).Error; err != nil {
		c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Program not found",
		})
	}
	if err := c.BodyParser(&program); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	var existingProgram models.Program
	// Check if the program already exists
	if err := initializers.DB.Where("program_name = ?", program.ProgramName).First(&existingProgram).Error; err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Program already exists",
		})
	}

	if err := initializers.DB.Save(&program).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not update program",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Program updated successfully",
		"program": program,
	})

}
