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
	// Get the uploaded file
	file, err := c.FormFile("file")
	if err != nil {
		fmt.Println("Error receiving file:", err) // Log to console for debugging
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
