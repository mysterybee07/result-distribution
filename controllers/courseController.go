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

	// Render the template with programs
	err := c.Render("dashboard/courses/addcourse", fiber.Map{
		"Programs": programs,
	})
	if err != nil {
		c.Status(fiber.StatusInternalServerError).SendString("Error rendering page")
		return err
	}

	return nil
}

func GetSemestersByProgramID(c *fiber.Ctx) error {
	programID := c.Params("programID")

	// Fetch semesters for the given programID
	var semesters []models.Semester
	if err := initializers.DB.Where("program_id = ?", programID).Find(&semesters).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch semesters",
		})
	}

	// Return semesters as JSON response
	return c.JSON(fiber.Map{
		"semesters": semesters,
	})
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

	// Validate and create courses
	for _, course := range payload.Courses {
		course.ProgramID = payload.ProgramID
		course.SemesterID = payload.SemesterID

		// Check if the course already exists for the same program
		var existingCourse models.Course
		if err := initializers.DB.Where("course_code = ? AND program_id = ?", course.CourseCode, payload.ProgramID).First(&existingCourse).Error; err == nil {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": fmt.Sprintf("Course '%s' already exists for the given program", course.Name),
			})
		}

		if err := initializers.DB.Create(&course).Error; err != nil {
			log.Printf("Error creating course: %v\n", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Could not create course",
			})
		}
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Courses created successfully",
	})
}

func GetFilteredCourses(c *fiber.Ctx) error {
	programID := c.Query("program_id")
	semesterID := c.Query("semester_id")

	var courses []models.Course
	if err := initializers.DB.Where("program_id = ? AND semester_id = ?", programID, semesterID).Find(&courses).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error fetching courses"})
	}

	return c.JSON(fiber.Map{"courses": courses})
}
