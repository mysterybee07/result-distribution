package controllers

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/mysterybee07/result-distribution-system/initializers"
	"github.com/mysterybee07/result-distribution-system/models"
)

func AddSemester(c *fiber.Ctx) error {
	// Fetch programs from the database
	var programs []models.Program
	if err := initializers.DB.Find(&programs).Error; err != nil {
		log.Printf("Failed to fetch programs: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to fetch programs")
	}

	// Render the add semester form with programs
	return c.Render("dashboard/semester/semester", fiber.Map{
		"Programs": programs,
	})
}

func StoreSemester(c *fiber.Ctx) error {
	semester := new(models.Semester)
	if err := c.BodyParser(semester); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	var program models.Program
	if err := initializers.DB.First(&program, semester.ProgramID).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Program not found",
		})
	}

	var existingSemester models.Semester
	// Check if a semester with the same name exists within the same program
	if err := initializers.DB.Where("name = ? AND program_id = ?", semester.Name, semester.ProgramID).First(&existingSemester).Error; err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Semester already exists for this program",
		})
	}

	if err := initializers.DB.Create(&semester).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not create semester",
		})
	}

	result := initializers.DB.Preload("Program").Find(&semester)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not retrieve semester with associations",
		})
	}

	return c.Redirect("/semesters")
}
