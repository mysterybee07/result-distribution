package utils

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/mysterybee07/result-distribution-system/initializers"
	"github.com/mysterybee07/result-distribution-system/models"
)

// Helper function to resolve college ID from either uint or string input
func ResolveCollegeID(collegeIDInput interface{}) (uint, error) {
	switch v := collegeIDInput.(type) {
	case float64:
		return uint(v), nil
	case string:
		var college models.College
		if err := initializers.DB.Where("college_name = ?", v).First(&college).Error; err != nil {
			return 0, fmt.Errorf("college not found for name: %s", v)
		}
		return college.ID, nil
	default:
		return 0, fmt.Errorf("invalid college identifier type")
	}
}

// Helper function to parse uint form values
func ParseUintFormValue(value string) (uint, error) {
	parsed, err := strconv.ParseUint(value, 10, 32)
	return uint(parsed), err
}

// Helper function to respond with an error
func RespondError(c *fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(fiber.Map{"error": message})
}
