package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mysterybee07/result-distribution-system/initializers"
	"github.com/mysterybee07/result-distribution-system/middleware/validation"
	"github.com/mysterybee07/result-distribution-system/models"
)

func Student(c *fiber.Ctx) error {
	var programs []models.Program
	if err := initializers.DB.Find(&programs).Error; err != nil {
		c.Status(fiber.StatusInternalServerError).SendString("Error fetching programs")
		return err
	}
	var batches []models.Batch
	if err := initializers.DB.Find(&batches).Error; err != nil {
		c.Status(fiber.StatusInternalServerError).SendString("Error fetching batches")
		return err
	}
	err := c.Render("dashboard/students/addstudent", fiber.Map{
		"Programs": programs,
		"Batches":  batches,
	})
	if err != nil {
		c.Status(fiber.StatusInternalServerError).SendString("Error rendering page")
		return err
	}
	return nil
}

func CreateStudents(c *fiber.Ctx) error {
	var input struct {
		BatchID   uint `json:"batch_id"`
		ProgramID uint `json:"program_id"`
		Students  []struct {
			Fullname           string `json:"fullname"`
			SymbolNumber       string `json:"symbol_number"`
			RegistrationNumber string `json:"registration_number"`
		} `json:"students"`
	}

	if err := c.BodyParser(&input); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "unable to parse request body",
		})
	}

	var students []models.Student
	for _, s := range input.Students {
		student := models.Student{
			SymbolNumber:       s.SymbolNumber,
			RegistrationNumber: s.RegistrationNumber,
			Fullname:           s.Fullname,
			BatchID:            input.BatchID,
			ProgramID:          input.ProgramID,
		}

		if err := validation.ValidateStudent(&student, false); err != nil {
			return c.Status(err.(*fiber.Error).Code).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		students = append(students, student)

	}
	if err := initializers.DB.Create(&students).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not add students",
		})
	}

	return c.JSON(fiber.Map{
		"message":  "Students added successfully",
		"students": students,
	})
}

func EditStudent(c *fiber.Ctx) error {
	err := c.Render("Students/edit", fiber.Map{})
	if err != nil {
		c.Status(fiber.StatusInternalServerError).SendString("Error rendering page")
		return err
	}
	return nil
}

func UpdateStudent(c *fiber.Ctx) error {
	id := c.Params("id") // Extract student ID from the URL parameter

	var input struct {
		Fullname           string `json:"fullname"`
		SymbolNumber       string `json:"symbol_number"`
		RegistrationNumber string `json:"registration_number"`
		BatchID            uint   `json:"batch_id"`
		ProgramID          uint   `json:"program_id"`
	}

	// Parse the JSON body
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Unable to parse request body",
		})
	}

	// Find the student by ID
	var student models.Student
	if err := initializers.DB.First(&student, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Student not found",
		})
	}

	// Update the fields
	student.Fullname = input.Fullname
	student.SymbolNumber = input.SymbolNumber
	student.RegistrationNumber = input.RegistrationNumber
	student.BatchID = input.BatchID
	student.ProgramID = input.ProgramID

	// Validate updated student data
	if err := validation.ValidateStudent(&student, true); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Save the updated student data
	if err := initializers.DB.Save(&student).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not update student",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Student updated successfully",
		"student": student,
	})
}

func GetStudents(c *fiber.Ctx) error {
	var students []models.Student

	// Use Preload to include Batch and Program associations
	result := initializers.DB.Preload("Batch").Preload("Program").Find(&students)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not retrieve students",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":  "Students retrieved successfully",
		"students": students,
	})
}

func GetStudentById(c *fiber.Ctx) error {
	id := c.Params("id")
	var student models.Student
	if err := initializers.DB.Preload("Batch").Preload("Program").Preload("Semester").First(&student, id).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Student not found",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Student retrieved successfully",
		"student": student,
	})
}

func DeleteStudent(c *fiber.Ctx) error {
	id := c.Params("id")
	var student models.Student
	if err := initializers.DB.First(&student, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Student not found",
		})
	}
	if err := initializers.DB.Delete(&student).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not delete student",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Student deleted successfully",
		"student": student,
	})
}

func GetFilteredStudents(c *fiber.Ctx) error {
	batchID := c.Query("batch_id")
	programID := c.Query("program_id")
	semesterID := c.Query("semester_id")

	var students []models.Student
	if err := initializers.DB.Preload("Batch").Preload("Program").Preload("Semester").
		Where("batch_id = ? AND program_id = ? AND current_semester = ?", batchID, programID, semesterID).
		Find(&students).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error fetching students",
		})
	}

	return c.JSON(fiber.Map{
		"students": students,
	})
}
