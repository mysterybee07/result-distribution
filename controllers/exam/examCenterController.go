package controller

import (
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/mysterybee07/result-distribution-system/utils"
)

func CentersAllocation(c *fiber.Ctx) error {
	// File paths for the colleges, centers, and preferences
	collegeFilePath := "data/schools_grade12_2081.tsv"
	centerFilePath := "data/centers_grade12_2081.tsv"
	preferenceFilePath := "data/prefs.tsv"

	// Parse colleges, centers, and preferences from the TSV files
	colleges, err := utils.ParseColleges(collegeFilePath)
	if err != nil {
		log.Printf("Error parsing colleges: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to parse colleges",
		})
	}

	centers, err := utils.ParseCenters(centerFilePath)
	if err != nil {
		log.Printf("Error parsing centers: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to parse centers",
		})
	}

	prefs, err := utils.ParsePreferences(preferenceFilePath)
	if err != nil {
		log.Printf("Error parsing preferences: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to parse preferences",
		})
	}

	// Perform the center assignment
	assignedCenters, unassignedColleges := utils.AssignCenters(colleges, centers, prefs)

	// Create files for assigned and unassigned data
	err = utils.WriteAssignedCentersToFile(assignedCenters, colleges, centers) // Pass both colleges and centers
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to write assigned centers",
		})
	}

	err = utils.WriteUnassignedCollegesToFile(unassignedColleges)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to write unassigned colleges",
		})
	}

	// Return success message
	return c.JSON(fiber.Map{
		"message": "Center assignments processed successfully!",
	})
}
