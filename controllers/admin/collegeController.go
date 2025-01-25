package controllers

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/mysterybee07/result-distribution-system/initializers"
	"github.com/mysterybee07/result-distribution-system/models"
	"github.com/mysterybee07/result-distribution-system/utils"
)

func UploadColleges(c *fiber.Ctx) error {
	// Check the Content-Type header to differentiate between file and JSON input
	contentType := c.Get("Content-Type")
	var college models.College

	if contentType == "application/json" {
		// Handle JSON input
		if err := c.BodyParser(&college); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid JSON payload",
				"err":   err.Error(),
			})
		}

		// Validate and save the college to the database
		if err := initializers.DB.Create(&college).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to add college", "details": err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"success": true,
			"college": college,
		})
	}

	// Handle file upload (default behavior)
	file, err := c.FormFile("file")
	if err != nil {
		fmt.Println("Error receiving file:", err) // Log to console for debugging
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "File upload failed",
		})
	}

	// Save the uploaded file to a temporary location
	filePath := fmt.Sprintf("./uploads/%s", file.Filename)
	err = c.SaveFile(file, filePath)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save file",
		})
	}

	// Call the ParseColleges function to parse the file
	colleges, err := utils.ParseColleges(filePath)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Save parsed colleges to the database
	for _, college := range colleges {
		if err := initializers.DB.FirstOrCreate(&college).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Failed to save colleges",
				"details": err.Error(),
			})
		}
	}

	return c.JSON(fiber.Map{
		"success": true, "colleges": colleges,
	})
}

func GetCenterCollegesByProgramAndBatch(c *fiber.Ctx) error {
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
	var results []struct {
		CollegeCode string `json:"college_code"`
		CollegeName string `json:"college_name"`
		Address     string `json:"address"`
		IsCenter    bool   `json:"is_center"`
		Capacity    int    `json:"capacity"`
	}

	if err := initializers.DB.Table("colleges").
		Select("colleges.college_code, colleges.college_name, colleges.address, COALESCE(capacity_and_counts.is_center, false) AS is_center, COALESCE(capacity_and_counts.capacity, 0) AS capacity").
		Joins("LEFT JOIN capacity_and_counts ON colleges.id = capacity_and_counts.college_id").
		Find(&results).Error; err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error": "Failed to fetch colleges",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"colleges": results,
	})
}

func AssignCenterAndCapacity(c *fiber.Ctx) error {
	// Check the content type of the request
	contentType := c.Get("Content-Type")

	if contentType == "application/json" {
		var requestData struct {
			BatchID   uint `json:"batch_id"`
			ProgramID uint `json:"program_id"`
			Records   []struct {
				CollegeName string `json:"college_name"`
				IsCenter    bool   `json:"is_center"`
				Capacity    int    `json:"capacity"`
			} `json:"records"`
		}

		// Parse JSON data from the request body
		if err := c.BodyParser(&requestData); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "failed to parse JSON data"})
		}

		// Process each record in the JSON array
		for _, record := range requestData.Records {
			if err := utils.ProcessRecord(record.CollegeName, requestData.BatchID, requestData.ProgramID, record.IsCenter, record.Capacity); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
			}
		}

	} else {
		// Assume data is coming from a TSV file
		file, err := os.Open("centers_and_capacities.tsv")
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to open file"})
		}
		defer file.Close()

		reader := csv.NewReader(file)
		reader.Comma = '\t' // Specify tab-delimited file
		records, err := reader.ReadAll()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to read file"})
		}

		for i, record := range records {
			if i == 0 {
				continue // Skip header row
			}

			collegeName := record[0]
			batchID, _ := strconv.ParseUint(record[1], 10, 32)
			programID, _ := strconv.ParseUint(record[2], 10, 32)
			isCenter, _ := strconv.ParseBool(record[3])
			capacity, _ := strconv.Atoi(record[4])

			if err := utils.ProcessRecord(collegeName, uint(batchID), uint(programID), isCenter, capacity); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
			}
		}
	}

	return c.JSON(fiber.Map{
		"message": "center status and capacity update completed successfully",
	})
}

func UpdateCollege(c *fiber.Ctx) error {
	// Parse the College ID from the URL parameter
	id := c.Params("id")

	// Retrieve the existing college record
	var college models.College
	if err := initializers.DB.First(&college, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "College not found", "details": err.Error(),
		})
	}

	// Parse the JSON payload to update the college
	var updateData models.College
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid JSON payload", "details": err.Error(),
		})
	}

	// Update the college record
	if err := initializers.DB.Model(&college).Updates(updateData).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update college", "details": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message":         "college updated successfully",
		"updated_college": college,
	})
}

func DeleteCollege(c *fiber.Ctx) error {
	// Parse the College ID from the URL parameter
	id := c.Params("id")

	// Check if the college exists
	var college models.College
	if err := initializers.DB.First(&college, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "College not found", "details": err.Error(),
		})
	}

	// Delete the college record
	if err := initializers.DB.Delete(&college).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete college", "details": err.Error(),
		})
	}

	return c.JSON(fiber.Map{

		"message": "College deleted successfully",
	})
}

func UpdateCapacity(c *fiber.Ctx) error {
	centerID := c.Params("id")

	// Parse request body to get the new capacity value
	type RequestBody struct {
		Capacity int `json:"capacity"`
	}
	var requestBody RequestBody

	if err := c.BodyParser(&requestBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	// Find the center by ID and ensure it's a center
	var center models.CapacityAndCount
	if err := initializers.DB.Where("id = ? AND is_center = ?", centerID, true).First(&center).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Center not found or not a valid center",
		})
	}

	// Update the capacity
	center.Capacity = requestBody.Capacity
	if err := initializers.DB.Save(&center).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update capacity",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":  "Capacity updated successfully",
		"center":   center.CollegeID,
		"capacity": center.Capacity,
	})
}
