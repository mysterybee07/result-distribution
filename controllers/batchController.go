package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mysterybee07/result-distribution-system/initializers"
	"github.com/mysterybee07/result-distribution-system/models"
)

func AddBatch(c *fiber.Ctx) error {
	err := c.Render("dashboard/batch/batch", fiber.Map{})
	if err != nil {
		c.Status(fiber.StatusInternalServerError).SendString("Error rendering page")
		return err
	}
	return nil
}

func CreateBatch(c *fiber.Ctx) error {
	// Initialize a new Batch instance
	batch := new(models.Batch)

	// Parse the request body into the Batch instance
	if err := c.BodyParser(batch); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	// Check if a batch with the same year already exists
	var existingBatch models.Batch
	if err := initializers.DB.Where("year = ?", batch.Year).First(&existingBatch).Error; err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Batch already exists",
		})
	}

	// Create the new batch in the database
	if err := initializers.DB.Create(batch).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not create batch",
		})
	}

	// Return a success response with the created batch
	return c.Redirect("/batches")
}
