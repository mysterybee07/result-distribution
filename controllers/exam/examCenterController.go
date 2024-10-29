package controller

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/mysterybee07/result-distribution-system/utils"
)

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
