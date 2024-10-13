package controllers

import (
	"os"
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
	if err := initializers.DB.Where("email = ? OR symbol_number = ?", loginData.Identifier, loginData.Identifier).First(&user).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "No user found for the email or symbol",
		})
	}

	// Check if the password is correct
	if !utils.CheckPasswordHash(loginData.Password, user.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Incorrect password or identifier",
		})
	}

	// Create JWT token
	token, err := utils.GenerateJwt(user.ID, user.Role, c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to Generate JWT tokens",
		})
	}

	// Return success response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User login successful",
		"user":    user,
		"token":   token, // Optional: you can return the token in the response as well
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
	// Update ImageURL if provided
	if updateData.ImageURL != "" {
		user.ImageURL = updateData.ImageURL
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

func GetLoginUser(c *fiber.Ctx) error {
	// Retrieve the JWT from the cookie
	cookie := c.Cookies("jwt")
	if cookie == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Missing or invalid JWT cookie",
		})
	}

	// Parse the JWT and extract the userID and role
	userID, role, err := utils.ParseJwt(cookie)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// Retrieve user information from the database based on userID
	var user models.User
	if err := initializers.DB.First(&user, userID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User not found",
		})
	}

	// Return the logged-in user's details
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"user": fiber.Map{
			"ID":    user.ID,
			"email": user.Email,
			"role":  role,
			// "name":  user.Name,
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
