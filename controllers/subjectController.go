package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mysterybee07/result-distribution-system/initializers"
	"github.com/mysterybee07/result-distribution-system/models"
)

func AddSubject(c *fiber.Ctx) error {
	err := c.Render("subjects/add", fiber.Map{})
	if err != nil {
		c.Status(fiber.StatusInternalServerError).SendString("Error rendering page")
		return err
	}
	return nil
}

func StoreSubject(c *fiber.Ctx) error {
	subject := new(models.Subject)

	// Parse the incoming JSON request body into the subject struct
	if err := c.BodyParser(subject); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	// Check if the program exists
	var program models.Program
	if err := initializers.DB.First(&program, subject.ProgramID).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Program not found",
		})
	}

	// Check if the semester exists
	var semester models.Semester
	if err := initializers.DB.First(&semester, subject.SemesterID).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Semester not found",
		})
	}

	// Check if the subject already exists for the same program (across any semester)
	var existingSubject models.Subject
	if err := initializers.DB.Where("name = ? AND program_id = ?", subject.Name, subject.ProgramID).First(&existingSubject).Error; err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Subject already exists for the given program",
		})
	}

	// Create the new subject
	if err := initializers.DB.Create(subject).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not create subject",
		})
	}

	// Preload the Program and Semester associations, including the nested Program within Semester
	if err := initializers.DB.Preload("Program").Preload("Semester.Program").First(&subject, subject.ID).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not retrieve subject with associations",
		})
	}

	// Return the created subject with a success message
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Subject created successfully",
		"subject": subject,
	})
}
