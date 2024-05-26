package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mysterybee07/result-distribution-system/initializers"
	"github.com/mysterybee07/result-distribution-system/middleware/validation"
	"github.com/mysterybee07/result-distribution-system/models"
)

func AddStudent(c *fiber.Ctx) error {
	err := c.Render("dashboard/students/addstudent", fiber.Map{})
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

	if err := validation.ValidateStudent(student, false); err != nil {
		return c.Status(err.(*fiber.Error).Code).JSON(fiber.Map{
			"error": err.Error(),
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
