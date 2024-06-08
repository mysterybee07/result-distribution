package controllers

import (
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/mysterybee07/result-distribution-system/initializers"
	"github.com/mysterybee07/result-distribution-system/middleware/validation"
	"github.com/mysterybee07/result-distribution-system/models"
)

func AddStudent(c *fiber.Ctx) error {
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

func StoreStudents(c *fiber.Ctx) error {
	// Parse JSON data
	var requestData struct {
		BatchID   uint             `json:"batch_id"`
		ProgramID uint             `json:"program_id"`
		Students  []models.Student `json:"students"`
	}

	if err := c.BodyParser(&requestData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	// Begin a transaction
	tx := initializers.DB.Begin()
	if tx.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not begin transaction",
		})
	}

	// Iterate over the students and validate/save each one
	for _, student := range requestData.Students {
		student.BatchID = requestData.BatchID
		student.ProgramID = requestData.ProgramID

		// Validate the student
		if err := validation.ValidateStudent(&student, false); err != nil {
			tx.Rollback()
			return c.Status(err.(*fiber.Error).Code).JSON(fiber.Map{
				"message": err.Error(),
			})
		}

		// Save the student
		if err := tx.Create(&student).Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Could not create student",
			})
		}
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not commit transaction",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Students created successfully",
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
	id := c.Params("id")
	var student models.Student

	if err := initializers.DB.First(&student, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Student not found",
		})
	}

	if err := c.BodyParser(&student); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	if err := validation.ValidateStudent(&student, true); err != nil {
		return c.Status(err.(*fiber.Error).Code).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Update the student
	if err := initializers.DB.Save(&student).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not update student",
		})
	}
	if err := initializers.DB.Preload("Batch").Preload("Program").First(&student, student.ID).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not retrieve student with associations",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
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

// example

func PostStudents(c *fiber.Ctx) error {
	var req struct {
		BatchID   string           `json:"batch_id"`
		ProgramID string           `json:"program_id"`
		Students  []models.Student `json:"students"`
	}

	if err := c.BodyParser(&req); err != nil {
		log.Printf("Error parsing body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Convert batch_id and program_id to uint
	batchID, err := strconv.ParseUint(req.BatchID, 10, 64)
	if err != nil {
		log.Printf("Error parsing batch_id: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid batch_id"})
	}

	programID, err := strconv.ParseUint(req.ProgramID, 10, 64)
	if err != nil {
		log.Printf("Error parsing program_id: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid program_id"})
	}

	// Create students
	for _, student := range req.Students {
		student.BatchID = uint(batchID)
		student.ProgramID = uint(programID)
		if err := initializers.DB.Create(&student).Error; err != nil {
			log.Printf("Error creating student: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Students created successfully"})
}
