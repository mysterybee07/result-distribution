package controllers

import (
	"log"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/mysterybee07/result-distribution-system/initializers"
	"github.com/mysterybee07/result-distribution-system/middleware/validation"
	"github.com/mysterybee07/result-distribution-system/models"
	"github.com/mysterybee07/result-distribution-system/utils"
)

// var initializers *gorm.DB
func Home(c *fiber.Ctx) error {
	err := c.Render("index", fiber.Map{})
	if err != nil {
		c.Status(fiber.StatusInternalServerError).SendString("Error rendering page")
		return err
	}
	return nil

}
func Register(c *fiber.Ctx) error {
	var batches []models.Batch
	if err := initializers.DB.Find(&batches).Error; err != nil {
		c.Status(fiber.StatusInternalServerError).SendString("Error fetching batches")
		return err
	}
	var programs []models.Program
	if err := initializers.DB.Find(&programs).Error; err != nil {
		c.Status(fiber.StatusInternalServerError).SendString("Error fetching programs")
		return err
	}
	err := c.Render("users/signup", fiber.Map{
		"Programs": programs,
		"Batches":  batches,
	})
	if err != nil {
		c.Status(fiber.StatusInternalServerError).SendString("Error rendering page")
		return err
	}
	return nil
}

// StoreRegister handles the registration of a new user

func StoreRegister(c *fiber.Ctx) error {
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "unable to parse json data",
		})
	}

	hashedPassword, err := models.HashPassword(user.Password)
	if err != nil {
		log.Println("Failed to hash password")
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Failed to process password",
		})
	}
	user.Password = hashedPassword

	if err := validation.ValidateUser(&user, false); err != nil {
		log.Println("Validation Failed")
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if err := initializers.DB.Create(&user).Error; err != nil {
		log.Println("Failed to create users")
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "User creation Failed",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User created successfully",
		"data":    user,
	})
}

func Login(c *fiber.Ctx) error {
	err := c.Render("users/login", fiber.Map{})
	if err != nil {
		c.Status(fiber.StatusInternalServerError).SendString("Error rendering page")
		return err
	}
	return nil
}

// LoginUser handles user login

func LoginUser(c *fiber.Ctx) error {
	type LoginData struct {
		Identifier string `json:"identifier"`
		Password   string `json:"password"`
	}

	var loginData LoginData

	// Parse the login data
	if err := c.BodyParser(&loginData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Unable to parse login data",
		})
	}

	// Find user by email or symbol
	var user models.User
	if err := initializers.DB.Where("email = ? OR symbol = ?", loginData.Identifier, loginData.Identifier).First(&user).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "No user found for the email or symbol",
		})
	}

	// Check if the password is correct
	if !models.CheckPasswordHash(loginData.Password, user.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Incorrect password or identifier",
		})
	}

	// Create JWT token
	_, err := utils.GenerateJwt(user.ID, user.Role, c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to Generate JWT tokens",
		})
	}

	// Return success response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User login successful",
		"user":    user,
		// "token":   token, // Optional: you can return the token in the response as well
	})
}

func UpdateUser(c *fiber.Ctx) error {
	userID := c.Params("id")

	// Fetch the existing user
	var existingUser models.User
	if err := initializers.DB.First(&existingUser, userID).Error; err != nil {
		log.Println("User not found:", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User not found",
		})
	}

	// Parse the form data
	form, err := c.MultipartForm()
	if err != nil {
		log.Println("Failed to get multipart form data:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to get multipart form data",
		})
	}

	// Update User instance with new data
	batchID, err := strconv.ParseUint(form.Value["batch_id"][0], 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid batch_id",
		})
	}
	batchIDPtr := uint(batchID)
	existingUser.BatchID = &batchIDPtr

	programID, err := strconv.ParseUint(form.Value["program_id"][0], 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid program_id",
		})
	}
	programIDPtr := uint(programID)
	existingUser.ProgramID = &programIDPtr

	existingUser.Symbol = form.Value["symbol"][0]
	existingUser.Registration = form.Value["registration"][0]
	existingUser.Email = form.Value["email"][0]
	existingUser.Password = form.Value["password"][0]
	// existingUser.Terms, err = strconv.ParseBool(form.Value["terms"][0])
	// if err != nil {
	// 	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
	// 		"message": "Invalid terms value",
	// 	})
	// }
	existingUser.Role = form.Value["role"][0]

	// Handle the image upload
	files := form.File["image_url"]
	if len(files) == 0 {
		log.Println("No image file found in the form data")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "No image file found in the form data",
		})
	}
	file := files[0]
	fileName := utils.RandLetter(5) + "-" + utils.SanitizeFileName(file.Filename)
	if len(files) > 0 {
		file := files[0] // assuming only one image file is uploaded
		fileName = utils.RandLetter(5) + "-" + utils.SanitizeFileName(file.Filename)
		filePath := filepath.Join("./static/images/uploads", fileName)

		log.Println("Saving file to:", filePath)
		if err := c.SaveFile(file, filePath); err != nil {
			log.Println("Failed to save image file:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to save image file",
			})
		}
	}

	// Set image URL regardless of upload success (empty string if no upload)
	existingUser.ImageURL = "/static/images/uploads/" + fileName

	// Validate the updated user data
	if err := validation.ValidateUser(&existingUser, true); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// Hash password if changed
	if len(existingUser.Password) < 8 {
		log.Println("Password too short")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Password must be at least 8 characters long",
		})
	} else if existingUser.Password != form.Value["password"][0] {
		hashedPassword, err := models.HashPassword(existingUser.Password)
		if err != nil {
			log.Println("Failed to hash password:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to hash password",
			})
		}
		existingUser.Password = hashedPassword
	}

	// Save the updated user in the database
	if err := initializers.DB.Save(&existingUser).Error; err != nil {
		log.Println("Failed to update user:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update user",
		})
	}

	// Return success message as JSON
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"user":    existingUser,
		"message": "Account updated successfully",
	})
}

// LogoutUser logs out the user by clearing the JWT cookie
func LogoutUser(c *fiber.Ctx) error {
	// Clear the cookie
	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Second), // Set the cookie to expire immediately
		HTTPOnly: true,
	})

	return c.Redirect("/login") // Redirect to login page or home page after logout
}

func ForgotPassword(c *fiber.Ctx) error {
	err := c.Render("users/forgot-password", fiber.Map{})
	if err != nil {
		c.Status(fiber.StatusInternalServerError).SendString("Error rendering page")
		return err
	}
	return nil
}
