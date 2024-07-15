package controllers

import (
	"errors"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/mysterybee07/result-distribution-system/initializers"
	"github.com/mysterybee07/result-distribution-system/models"
	"gorm.io/gorm"
)

func AddCourse(c *fiber.Ctx) error {
	// Fetch programs with their associated semesters
	var programs []models.Program
	if err := initializers.DB.Preload("Semesters").Find(&programs).Error; err != nil {
		c.Status(fiber.StatusInternalServerError).SendString("Error fetching programs")
		return err
	}

	// Fetch batches
	var batches []models.Batch
	if err := initializers.DB.Find(&batches).Error; err != nil {
		c.Status(fiber.StatusInternalServerError).SendString("Error fetching batches")
		return err
	}

	// Render the template with both programs and batches
	err := c.Render("dashboard/courses/addcourse", fiber.Map{
		"Programs": programs,
		"Batches":  batches,
	})
	if err != nil {
		c.Status(fiber.StatusInternalServerError).SendString("Error rendering page")
		return err
	}

	return nil
}

func GetSemestersByProgram(c *fiber.Ctx) error {
	programId := c.Params("id")
	var semesters []models.Semester
	if err := initializers.DB.Where("program_id = ?", programId).Find(&semesters).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error fetching semesters"})
	}
	return c.JSON(fiber.Map{"semesters": semesters})
}

// StoreCourse handles storing multiple courses in a single request
func StoreCourse(c *fiber.Ctx) error {
	// Parse incoming JSON request body into payload struct
	var payload models.CoursesPayload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid JSON",
		})
	}

	// Check if the program exists
	var program models.Program
	if err := initializers.DB.First(&program, payload.ProgramID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Program not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error",
		})
	}

	// Check if the semester exists
	var semester models.Semester
	if err := initializers.DB.First(&semester, payload.SemesterID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Semester not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error",
		})
	}

	// Create and validate courses
	var createdCourses []models.Course
	for _, course := range payload.Courses {
		// Check if the course already exists for the same program
		var existingCourse models.Course
		if err := initializers.DB.Where("name = ? AND program_id = ?", course.Name, payload.ProgramID).First(&existingCourse).Error; err == nil {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": fmt.Sprintf("Course '%s' already exists for the given program", course.Name),
			})
		}
		newCourse := models.Course{
			CourseCode:          course.CourseCode,
			Name:                course.Name,
			SemesterPassMarks:   course.SemesterPassMarks,
			PracticalPassMarks:  course.PracticalPassMarks,
			AssistantPassMarks:  course.AssistantPassMarks,
			SemesterTotalMarks:  course.SemesterTotalMarks,
			PracticalTotalMarks: course.PracticalTotalMarks,
			AssistantTotalMarks: course.AssistantTotalMarks,
			ProgramID:           payload.ProgramID,
			SemesterID:          payload.SemesterID,
		}

		if err := initializers.DB.Create(&newCourse).Error; err != nil {
			log.Printf("Error creating course: %v\n", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Could not create course",
			})
		}

		// Append created course to response
		createdCourses = append(createdCourses, newCourse)
	}

	// Return success response with created courses
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Courses created successfully",
		"courses": createdCourses,
	})
}
