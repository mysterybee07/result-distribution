package controller

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/mysterybee07/result-distribution-system/utils"
)

func UploadColleges(c *fiber.Ctx) error {
	// Get the uploaded file
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "File upload failed"})
	}

	// Save the uploaded file to a temporary location
	filePath := fmt.Sprintf("./uploads/%s", file.Filename)
	err = c.SaveFile(file, filePath)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save file"})
	}

	// Call the ParseColleges function (without batchID and programID)
	colleges, err := utils.ParseColleges(filePath)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true, "colleges": colleges})
}

func AssignCentersHandler(c *fiber.Ctx) error {
	// Parse batch and program IDs from the query parameters
	batchID, err := strconv.ParseUint(c.Query("batch_id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid batch_id",
		})
	}

	programID, err := strconv.ParseUint(c.Query("program_id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid program_id",
		})
	}

	// Call the AssignCenters function with the parsed batchID and programID
	assignments, err := utils.AssignCenters(uint(batchID), uint(programID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to assign centers: %s", err.Error()),
		})
	}

	// Write the results to a CSV file
	if err := utils.WriteResultToFile(assignments); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to write result to file: %s", err.Error()),
		})
	}

	// Return the assignments in the response
	return c.Status(fiber.StatusOK).JSON(assignments)
}
