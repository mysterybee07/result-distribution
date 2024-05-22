package controllers

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/mysterybee07/result-distribution-system/initializers"
	"github.com/mysterybee07/result-distribution-system/models"
)

func AddStudent(c *fiber.Ctx) error {
	err := c.Render("students/add", fiber.Map{})
	if err != nil {
		c.Status(fiber.StatusInternalServerError).SendString("Error rendering page")
		return err
	}
	return nil
}

func StoreStudent(c *fiber.Ctx) error {
	student := new(models.Student)
	if err := c.BodyParser(student); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	// Check if Batch and Program exist
	var batch models.Batch
	var program models.Program

	if err := initializers.DB.First(&batch, student.BatchID).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Batch not found",
		})
	}

	if err := initializers.DB.First(&program, student.ProgramID).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Program not found",
		})
	}

	//check if symbol number and registration number already exists
	var existingSymbol models.Student

	if err := initializers.DB.Where("symbol_number = ? AND batch_id = ? AND program_id = ?", student.SymbolNumber, student.BatchID, student.ProgramID).First(&existingSymbol).Error; err == nil {
		log.Println("Symbol Number already taken in users for the specified batch and program:", student.SymbolNumber)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Symbol Number is already taken for the specified batch and program",
		})
	}
	if err := initializers.DB.Where("registration = ? AND batch_id = ? AND program_id = ?", student.Registration, student.BatchID, student.ProgramID).First(&existingSymbol).Error; err == nil {
		log.Println("Registration Number already taken in users for the specified batch and program:", student.Registration)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Registration Number is already taken for the specified batch and program",
		})
	}

	// Save the student
	if err := initializers.DB.Create(&student).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not create student",
		})
	}
	if err := initializers.DB.Preload("Batch").Preload("Program").First(&student, student.ID).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not retrieve student with associations",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Student created successfully",
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
