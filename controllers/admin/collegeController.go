package controllers

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/mysterybee07/result-distribution-system/initializers"
	"github.com/mysterybee07/result-distribution-system/models"
	"github.com/mysterybee07/result-distribution-system/utils"
	"gorm.io/gorm"
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

func AssignCenterAndCapacity(c *fiber.Ctx) error {
	// Define a struct to hold the incoming JSON data
	type RequestBody struct {
		CollegeID uint `json:"college_id"`
		BatchID   uint `json:"batch_id"`
		ProgramID uint `json:"program_id"`
		IsCenter  bool `json:"is_center"`
		Capacity  int  `json:"capacity"`
	}

	// Parse JSON data from the request body
	var body RequestBody
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to parse request body"})
	}

	// Check if the college exists in the College model
	var college models.College
	if err := initializers.DB.First(&college, body.CollegeID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "College not found"})
	}

	// Try to find the existing CapacityAndCount record for the specified college, batch, and program
	var capacityAndCount models.CapacityAndCount
	result := initializers.DB.Where("college_id = ? AND batch_id = ? AND program_id = ?", body.CollegeID, body.BatchID, body.ProgramID).First(&capacityAndCount)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// If no record is found, create a new CapacityAndCount record with students_count set to 0
			capacityAndCount = models.CapacityAndCount{
				CollegeID:     body.CollegeID,
				BatchID:       body.BatchID,
				ProgramID:     body.ProgramID,
				StudentsCount: 0,
				IsCenter:      body.IsCenter,
				Capacity:      body.Capacity,
			}
			if err := initializers.DB.Create(&capacityAndCount).Error; err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create new record"})
			}
		} else {
			// If there's another error, return it
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve record"})
		}
	} else {
		// If the record exists, update the center status and capacity
		capacityAndCount.IsCenter = body.IsCenter
		capacityAndCount.Capacity = body.Capacity

		if err := initializers.DB.Save(&capacityAndCount).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update record"})
		}
	}

	return c.JSON(fiber.Map{
		"message":            "Center status and capacity updated successfully",
		"capacity_and_count": capacityAndCount,
	})
}
