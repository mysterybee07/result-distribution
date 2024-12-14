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

func Course(c *fiber.Ctx) error {
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

// StoreCourse handles storing multiple courses in a single request
func CreateCourses(c *fiber.Ctx) error {
	var payload models.CoursesPayload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request body",
		})
	}

	// Check if courses were provided
	if len(payload.Courses) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No courses provided. Please specify at least one course.",
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

	var courses []models.Course
	// Validate and create courses
	for _, course := range payload.Courses {
		// Validate that course has necessary fields
		if course.CourseCode == "" || course.Name == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Course code and name cannot be empty",
			})
		}

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
		courses = append(courses, course)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Courses created successfully",
		"courses": courses,
	})
}

func UpdateCourse(c *fiber.Ctx) error {
	id := c.Params("id")

	var course models.Course

	if err := initializers.DB.First(&course, id).Error; err != nil {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"error": fmt.Sprintf("Course with id:%s is not found", id),
		})
	}

	if err := c.BodyParser(&course); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"error": "Failed to parse request body",
		})
	}

	var existingCourse models.Course
	if err := initializers.DB.Where("course_code = ? AND program_id = ?", course.CourseCode, course.ProgramID).First(&existingCourse).Error; err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": fmt.Sprintf("Course '%s' already exists for the given program", course.Name),
		})
	}

	if err := initializers.DB.Save(&course).Error; err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error": "Failed to update course",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Courses updated successfully",
		"course":  course,
	})
}

func GetFilteredCourses(c *fiber.Ctx) error {
	programID := c.Query("program_id")
	semesterID := c.Query("semester_id")

	var courses []models.Course
	if err := initializers.DB.Where("program_id = ? AND semester_id = ?", programID, semesterID).Find(&courses).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error fetching courses"})
	}

	return c.JSON(fiber.Map{
		"courses": courses,
	})
}

func GetCourseById(c *fiber.Ctx) error {
	id := c.Params("id")
	var course models.Course
	if err := initializers.DB.First(&course, id).Error; err != nil {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"error": "Course not found",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Course retrieve successfully",
		"course":  course,
	})
}
