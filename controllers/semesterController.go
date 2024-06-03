package controllers

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/mysterybee07/result-distribution-system/initializers"
	"github.com/mysterybee07/result-distribution-system/models"
)

func AddSemester(c *fiber.Ctx) error {
	var programs []models.Program
	if err := initializers.DB.Preload("Semesters").Find(&programs).Error; err != nil {
		log.Printf("Failed to fetch programs: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to fetch programs")
	}

	err := c.Render("dashboard/semester/semester", fiber.Map{
		"Programs": programs,
	})
	if err != nil {
		log.Printf("Failed to render page: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Error rendering page")
	}

	return nil
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

func EditSemester(c *fiber.Ctx) error {
	err := c.Render("dashboard/semester/editSemester", fiber.Map{})
	if err != nil {
		c.Status(fiber.StatusInternalServerError).SendString("Error rendering page")
		return err
	}
	return nil
}

func UpdateSemester(c *fiber.Ctx) error {
	id := c.Params("id")

	var semester models.Semester

	if err := initializers.DB.Find(&semester, id).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Semester not found")
	}

	if err := c.BodyParser(&semester); err != nil {
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

	if err := initializers.DB.Save(&semester).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not update semester",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":  "Semester updated successfully",
		"semester": semester,
	})
}
