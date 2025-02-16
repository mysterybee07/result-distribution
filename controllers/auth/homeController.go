package controllers

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/mysterybee07/result-distribution-system/initializers"
	"github.com/mysterybee07/result-distribution-system/middleware/validation"
	"github.com/mysterybee07/result-distribution-system/models"
	"github.com/mysterybee07/result-distribution-system/utils"
	"gorm.io/gorm"
)

// StoreRegister handles the registration of a new user

func StoreRegister(c *fiber.Ctx) error {
	// Create a new user input instance
	var userInput models.UserInput
	if err := c.BodyParser(&userInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request payload",
		})
	}

	// Create a new user instance
	user := models.User{
		BatchID:            &userInput.BatchID,
		ProgramID:          &userInput.ProgramID,
		SymbolNumber:       userInput.SymbolNumber,
		RegistrationNumber: userInput.RegistrationNumber,
		Email:              userInput.Email,
		Password:           userInput.Password,
		Role:               userInput.Role,
	}

	// Hash the password
	hashedPassword, err := utils.HashPassword(userInput.Password)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to process password",
		})
	}
	user.Password = hashedPassword

	// Validate user data
	if err := validation.ValidateUser(&user, false); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// Save the user to the database first (before saving the image to the file system)
	if err := initializers.DB.Create(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "User creation failed",
			"error":   err.Error(),
		})
	}

	// Handle image upload and get the file path
	imageURL, err := utils.UploadImage(c)
	if err != nil {
		// If image upload fails, rollback the user creation and delete the user record from DB
		initializers.DB.Delete(&user)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error uploading image: " + err.Error(),
		})
	}

	// Now that both the user is created and image is uploaded, update the user with the image URL
	user.ImageURL = imageURL
	if err := initializers.DB.Save(&user).Error; err != nil {
		// In case of error while updating user with image, delete the image file and rollback the user creation
		os.Remove(imageURL) // Delete the uploaded image from the folder
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update user with image",
			"error":   err.Error(),
		})
	}

	// Return a success message
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User created successfully",
		"user":    user,
	})
}

func LoginUser(c *fiber.Ctx) error {
	type LoginData struct {
		Identifier string `json:"identifier"`
		Password   string `json:"password"`
	}

	var loginData LoginData

	// Parse the login data
	if err := c.BodyParser(&loginData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": fiber.Map{
				"message": "Unable to parse login data",
			},
		})
	}

	var user models.User

	// Find user by email or symbol
	if err := initializers.DB.Where("email = ? OR symbol_number = ?", loginData.Identifier, loginData.Identifier).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"errors": fiber.Map{
					"identifier": "User not found for the provided email or symbol",
				},
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errors": fiber.Map{
				"message": "Database error occurred",
			},
		})
	}

	// Check if the password is correct
	if !utils.CheckPasswordHash(loginData.Password, user.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"errors": fiber.Map{
				"password": "Incorrect password or identifier",
			},
		})
	}

	// Generate JWT token after successful login
	_, err := utils.GenerateJwt(user.ID, user.Role, c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errors": fiber.Map{"message": "Failed to generate JWT token"},
		})
	}

	// Create a response DTO
	userResponse := struct {
		ID    uint   `json:"id"`
		Email string `json:"email"`
		Role  string `json:"role"`
	}{
		ID:    user.ID,
		Email: user.Email,
		Role:  user.Role,
	}

	// Return success response with user details and token
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":    fiber.StatusOK,
		"message": "User login successful",
		"user":    userResponse,
		// "token":   token,
	})
}

func UpdateUser(c *fiber.Ctx) error {
	// Take id
	id := c.Params("id")
	var user models.User

	// Find user
	if err := initializers.DB.First(&user, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User with the UserID not found",
		})
	}

	var updateData models.User
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to parse request body",
		})
	}

	if err := validation.ValidateUser(&updateData, true); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	// Update the user fields only if they are provided
	if updateData.Email != "" && updateData.Email != user.Email {
		user.Email = updateData.Email
	}

	if updateData.Password != "" {
		hashpassword, err := utils.HashPassword(updateData.Password)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to hash password",
			})
		}
		user.Password = hashpassword
	}

	// Handle image upload using the UpdateImage function
	newImagePath, err := utils.UpdateImage(c, user.ImageURL)
	if err != nil && err.Error() != "no image file found in the form data" {
		// If there was an error other than no image being uploaded, return the error
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update image: " + err.Error(),
		})
	}

	// If a new image was uploaded, update the ImageURL field
	if newImagePath != "" {
		user.ImageURL = newImagePath
	}

	if err := initializers.DB.Save(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update user in the database",
		})
	}

	// Return success message as JSON
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		// "user":    existingUser,
		"message": "Account updated successfully",
	})
}

// LogoutUser logs out the user by clearing the JWT cookie
func LogoutUser(c *fiber.Ctx) error {
	// Clear the cookie
	cookies := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Second), // Set the cookie to expire immediately
		HTTPOnly: true,
	}
	c.Cookie(&cookies)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "logout successfully",
	})
}

func AuthorizedUser(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string) // Retrieve user ID from locals

	// Retrieve user information from database
	var user models.User
	if err := initializers.DB.First(&user, userID).Error; err != nil {
		log.Println("User not found in database")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Return logged-in user's details
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": fiber.Map{
			"ID":    user.ID,
			"email": user.Email,
			"role":  user.Role, // Make sure `role` is set correctly
		},
		"message": "User data retrieved successfully",
	})
}

func GetAllUsers(c *fiber.Ctx) error {
	var users []models.User

	if err := initializers.DB.Find(&users).Error; err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Unable to retrive users",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "user retrieved successfully",
		"users":   users,
	})
}

func GetUserById(c *fiber.Ctx) error {
	id := c.Params("id")

	var user models.User

	if err := initializers.DB.First(&user, id).Error; err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "User with id not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"user": user,
	})
}

func ForgotPassword(c *fiber.Ctx) error {
	err := c.Render("users/forgot-password", fiber.Map{})
	if err != nil {
		c.Status(fiber.StatusInternalServerError).SendString("Error rendering page")
		return err
	}
	return nil
}
