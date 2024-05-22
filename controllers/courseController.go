package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mysterybee07/result-distribution-system/initializers"
	"github.com/mysterybee07/result-distribution-system/models"
)

func AddCourse(c *fiber.Ctx) error {
	err := c.Render("subjects/add", fiber.Map{})
	if err != nil {
		c.Status(fiber.StatusInternalServerError).SendString("Error rendering page")
		return err
	}
	return nil
}

func StoreCourse(c *fiber.Ctx) error {
	course := new(models.Course)

	// Parse the incoming JSON request body into the course struct
	if err := c.BodyParser(course); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	// Check if the program exists
	var program models.Program
	if err := initializers.DB.First(&program, course.ProgramID).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Program not found",
		})
	}

	// Check if the semester exists
	var semester models.Semester
	if err := initializers.DB.First(&semester, course.SemesterID).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Semester not found",
		})
	}

	// Check if the course already exists for the same program (across any semester)
	var existingCourse models.Course
	if err := initializers.DB.Where("name = ? AND program_id = ?", course.Name, course.ProgramID).First(&existingCourse).Error; err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "course already exists for the given program",
		})
	}

	// Create the new course
	if err := initializers.DB.Create(course).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not create course",
		})
	}

	// Preload the Program and Semester associations, including the nested Program within Semester
	if err := initializers.DB.Preload("Program").Preload("Semester.Program").First(&course, course.ID).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not retrieve course with associations",
		})
	}

	// Return the created course with a success message
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "course created successfully",
		"course":  course,
	})
}
