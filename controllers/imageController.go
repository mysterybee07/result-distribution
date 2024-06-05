package controllers

import (
	"log"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/mysterybee07/result-distribution-system/utils"
)

func Upload(c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		log.Println("Failed to get multipart form data:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to get multipart form data",
		})
	}

	files := form.File["image"]
	fileName := ""
	if len(files) > 0 {
		file := files[0] // assuming only one image file is uploaded
		fileName := utils.RandLetter(5) + "-" + utils.SanitizeFileName(file.Filename)
		filePath := filepath.Join("./static/images", fileName)

		if err := c.SaveFile(file, filePath); err != nil {
			log.Println("Failed to save image file:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to save image file",
			})
		}

	}
	return c.JSON(fiber.Map{
		"url": "http://localhost:8000/static/images" + fileName, // Adjust based on your server config

	})
}
