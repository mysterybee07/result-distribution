package controllers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/mysterybee07/result-distribution-system/initializers"
	"github.com/mysterybee07/result-distribution-system/models"
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

func GetCenterColleges(c *fiber.Ctx) error {
	batchID := c.Query("batch_id")
	programID := c.Query("program_id")

	type CenterResponse struct {
		CollegeName   string `json:"college_name"`
		Address       string `json:"address"`
		Capacity      int    `json:"capacity"`
		StudentsCount int    `json:"students_count"`
	}

	var centers []CenterResponse

	// Join College and select required fields
	if err := initializers.DB.Model(&models.CapacityAndCount{}).
		Joins("JOIN colleges ON colleges.id = capacity_and_counts.college_id").
		Where("batch_id = ? AND program_id = ? AND is_center = ?", batchID, programID, false).
		Select("colleges.college_name as college_name,colleges.address, capacity_and_counts.capacity, capacity_and_counts.students_count").
		Scan(&centers).Error; err != nil {

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch center colleges",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"centers": centers,
	})
}

func GetColleges(c *fiber.Ctx) error {
	var colleges []models.College

	if err := initializers.DB.Find(&colleges).Error; err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error": "College not found",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"center": colleges,
	})
}

// func AssignCenterAndCapacity(c *fiber.Ctx) error{
// 	var
// 	collegeID := c.Params("id")
// 	batchID := c.Params("id")
// 	programID := c.Params("id")

// }
