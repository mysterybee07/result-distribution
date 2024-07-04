package controllers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/mysterybee07/result-distribution-system/initializers"
	"github.com/mysterybee07/result-distribution-system/models"
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

func StoreCourse(c *fiber.Ctx) error {
	// Struct to parse the incoming JSON request body
	type CourseInput struct {
		CourseCode          string `json:"course_code"`
		Name                string `json:"name" validate:"required"`
		SemesterPassMarks   int    `json:"semester_pass_marks" validate:"required"`
		PracticalPassMarks  *int   `json:"practical_pass_marks,omitempty"`
		AssistantPassMarks  *int   `json:"assistant_pass_marks,omitempty"`
		SemesterTotalMarks  int    `json:"semester_total_marks" validate:"required"`
		PracticalTotalMarks *int   `json:"practical_total_marks,omitempty"`
		AssistantTotalMarks *int   `json:"assistant_total_marks,omitempty"`
	}

	type RequestPayload struct {
		ProgramID  uint          `json:"program_id" validate:"required"`
		SemesterID uint          `json:"semester_id" validate:"required"`
		Courses    []CourseInput `json:"courses" validate:"required"`
	}

	payload := new(RequestPayload)

	// Parse the incoming JSON request body into the payload struct
	if err := c.BodyParser(payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	// Check if the program exists
	var program models.Program
	if err := initializers.DB.First(&program, payload.ProgramID).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Program not found",
		})
	}

	// Check if the semester exists
	var semester models.Semester
	if err := initializers.DB.First(&semester, payload.SemesterID).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Semester not found",
		})
	}

	// Create and validate courses
	var createdCourses []models.Course
	for _, courseInput := range payload.Courses {
		// Check if the course already exists for the same program
		var existingCourse models.Course
		if err := initializers.DB.Where("name = ? AND program_id = ?", courseInput.Name, payload.ProgramID).First(&existingCourse).Error; err == nil {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": fmt.Sprintf("Course '%s' already exists for the given program", courseInput.Name),
			})
		}

		// Create the new course
		course := models.Course{
			CourseCode:          courseInput.CourseCode,
			Name:                courseInput.Name,
			SemesterPassMarks:   courseInput.SemesterPassMarks,
			PracticalPassMarks:  courseInput.PracticalPassMarks,
			AssistantPassMarks:  courseInput.AssistantPassMarks,
			SemesterTotalMarks:  courseInput.SemesterTotalMarks,
			PracticalTotalMarks: courseInput.PracticalTotalMarks,
			AssistantTotalMarks: courseInput.AssistantTotalMarks,
			ProgramID:           payload.ProgramID,
			SemesterID:          payload.SemesterID,
		}

		if err := initializers.DB.Create(&course).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Could not create course",
			})
		}

		// Preload the Program and Semester associations
		if err := initializers.DB.Preload("Program").Preload("Semester.Program").First(&course, course.ID).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Could not retrieve course with associations",
			})
		}

		createdCourses = append(createdCourses, course)
	}

	// Return the created courses with a success message
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Courses created successfully",
		"courses": createdCourses,
	})
}
