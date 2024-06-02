package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mysterybee07/result-distribution-system/initializers"
	"github.com/mysterybee07/result-distribution-system/models"
)

func AddBatch(c *fiber.Ctx) error {
	var batch []models.Batch
	if err := initializers.DB.Find(&batch).Error; err != nil {
		c.Status(fiber.StatusInternalServerError).SendString("Error fetching batch")
		return err
	}

	err := c.Render("dashboard/batch/batch", fiber.Map{
		"Batches": batch,
	})
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

func EditBatch(c *fiber.Ctx) error {
	err := c.Render("/batch/editbatch", fiber.Map{})
	if err != nil {
		c.Status(fiber.StatusInternalServerError).SendString("Error rendering page")
		return err
	}
	return nil
}

func UpdateBatch(c *fiber.Ctx) error {
	id := c.Params("id")

	var batch models.Batch
	if err := initializers.DB.First(&batch, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Batch not found",
		})
	}
	if err := c.BodyParser(&batch); err != nil {
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
	if err := initializers.DB.Save(&batch).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not update batch",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Batch updated successfully",
		"batch":   batch,
	})
}
