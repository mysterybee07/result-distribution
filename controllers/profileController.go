package controllers

import (
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/mysterybee07/result-distribution-system/initializers"
	"github.com/mysterybee07/result-distribution-system/models"
)

func Profile(c *fiber.Ctx) error {
	err := c.Render("users/profile", fiber.Map{})
	if err != nil {
		c.Status(fiber.StatusInternalServerError).SendString("Error rendering page")
		return err
	}
	return nil
}
func GetUserProfile(c *fiber.Ctx) error {
	// Get userID from context
	userID, ok := c.Locals("userID").(string)
	if !ok {
		log.Println("userID not found in context")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	// Convert userID to int
	id, err := strconv.Atoi(userID)
	if err != nil {
		log.Printf("Failed to convert userID to integer: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Server error"})
	}

	// Fetch user from database
	var user models.User
	if err := initializers.DB.First(&user, id).Error; err != nil {
		log.Printf("User not found: %v\n", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	// Fetch student information
	var student models.Student
	if err := initializers.DB.Where("symbol_number = ?", user.Symbol).Preload("Batch").Preload("Program").First(&student).Error; err != nil {
		log.Printf("Student not found for user with email %s: %v\n", user.Email, err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Student not found"})
	}

	// Fetch marks for the student
	var marks []models.Mark
	if err := initializers.DB.Where("student_id = ?", student.ID).Preload("Course").Find(&marks).Error; err != nil {
		log.Printf("Failed to find marks for student with ID %d: %v\n", student.ID, err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Marks not found"})
	}

	// Return user profile
	return c.Render("users/profile", fiber.Map{
		"user":    user,
		"student": student,
		"marks":   marks,
	})
}
