package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mysterybee07/result-distribution-system/initializers"
	"github.com/mysterybee07/result-distribution-system/models"
)

func AddBatch(c *fiber.Ctx) error {
	err := c.Render("dashboard/batch/addbatch", fiber.Map{})
	if err != nil {
		c.Status(fiber.StatusInternalServerError).SendString("Error rendering page")
		return err
	}
	return nil
}

func CreateBatch(c *fiber.Ctx) error {
	batch := new(models.Batch)
	if err := c.BodyParser(&batch); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}
	var existingBatch models.Batch

	if err := initializers.DB.Where("year=?", &batch.Year).First(&existingBatch).Error; err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Batch already exists",
		})
	}

	if err := initializers.DB.Create(&batch).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not create batch",
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Batch created successfully",
		"batch":   batch,
	})
}
