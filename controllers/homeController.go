package controllers

import (
	"log"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/mysterybee07/result-distribution-system/initializers"
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
	err := c.Render("users/signup", fiber.Map{})
	if err != nil {
		c.Status(fiber.StatusInternalServerError).SendString("Error rendering page")
		return err
	}
	return nil
}

// StoreRegister handles the registration of a new user

func StoreRegister(c *fiber.Ctx) error {
	// Parse the form data
	form, err := c.MultipartForm()
	if err != nil {
		log.Println("Failed to get multipart form data:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to get multipart form data",
		})
	}

	// Create a new User instance
	var data models.User

	// Parse form values into the User struct
	batchID, err := strconv.ParseUint(form.Value["batch_id"][0], 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid batch_id",
		})
	}
	batchIDPtr := uint(batchID)
	data.BatchID = &batchIDPtr

	programID, err := strconv.ParseUint(form.Value["program_id"][0], 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid program_id",
		})
	}
	programIDPtr := uint(programID)
	data.ProgramID = &programIDPtr

	data.Symbol = form.Value["symbol"][0]
	data.Registration = form.Value["registration"][0]
	data.Email = form.Value["email"][0]
	data.Password = form.Value["password"][0]
	data.Terms, err = strconv.ParseBool(form.Value["terms"][0])
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid terms value",
		})
	}
	data.Role = form.Value["role"][0]

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
	data.ImageURL = "/static/images/uploads/" + fileName

	// Check if the user is an admin
	if data.Role == "admin" || data.Role == "superadmin" {
		// Validate required fields for admin
		if data.Email == "" || data.Password == "" || data.Symbol == "" {
			log.Println("Missing required fields for admin")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Email, Password, and Symbol are required for admin",
			})
		}
		data.BatchID = nil
		data.ProgramID = nil
	} else {
		// Validate required fields for regular user
		if data.BatchID == nil || data.ProgramID == nil || data.Symbol == "" || data.Registration == "" || data.Email == "" || data.Password == "" {
			log.Println("Missing required fields for user")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "All fields are required for user",
			})
		}

		// Check if symbol and registration exist in the students table for the given batch and program
		var student models.Student
		if err := initializers.DB.Where("symbol_number = ? AND registration = ? AND batch_id = ? AND program_id = ?",
			data.Symbol, data.Registration, *data.BatchID, *data.ProgramID).First(&student).Error; err != nil {
			log.Println("Student record not found:", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Invalid symbol or registration for the specified batch and program",
			})
		}
	}

	// Check if the password is at least 8 characters long
	if len(data.Password) < 8 {
		log.Println("Password too short")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Password must be at least 8 characters long",
		})
	}

	// Validate email
	if !utils.ValidateEmail(data.Email) {
		log.Println("Invalid email format")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid email",
		})
	}

	// Check if email already exists
	var existingUser models.User
	if err := initializers.DB.Where("email = ?", data.Email).First(&existingUser).Error; err == nil {
		log.Println("Email already taken:", data.Email)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Email is already taken",
		})
	}

	// Check if symbol number and registration number are unique in users table
	if err := initializers.DB.Where("symbol = ? AND batch_id = ? AND program_id = ?", data.Symbol, data.BatchID, data.ProgramID).First(&existingUser).Error; err == nil {
		log.Println("Symbol Number already taken in users for the specified batch and program:", data.Symbol)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Symbol Number is already taken for the specified batch and program",
		})
	}
	if err := initializers.DB.Where("registration = ? AND batch_id = ? AND program_id = ?", data.Registration, data.BatchID, data.ProgramID).First(&existingUser).Error; err == nil {
		log.Println("Registration Number already taken in users for the specified batch and program:", data.Registration)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Registration Number is already taken for the specified batch and program",
		})
	}

	// Hash password
	hashedPassword, err := models.HashPassword(data.Password)
	if err != nil {
		log.Println("Failed to hash password:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to hash password",
		})
	}
	data.Password = hashedPassword

	// Create user in database
	if err := initializers.DB.Create(&data).Error; err != nil {
		log.Println("Failed to create user:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user",
		})
	}

	// Return success message as JSON
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"user":    data,
		"message": "Account created successfully",
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
	// Define a struct to parse the login form data
	type LoginData struct {
		Identifier string `json:"identifier" form:"identifier"`
		Password   string `json:"password" form:"password"`
	}

	// Parse the login form data
	var loginData LoginData
	if err := c.BodyParser(&loginData); err != nil {
		log.Println("Unable to parse form data:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid form data",
		})
	}

	// Find the user by email or symbol number
	var user models.User
	if err := initializers.DB.Where("email = ? OR symbol = ?", loginData.Identifier, loginData.Identifier).First(&user).Error; err != nil {
		log.Printf("User not found with identifier: %s\n", loginData.Identifier)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid identifier or password",
		})
	}

	// Verify the provided password against the stored hashed password
	if !models.CheckPasswordHash(loginData.Password, user.Password) {
		log.Println("Invalid password")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid identifier or password",
		})
	}

	// Generate JWT token with user ID and role
	token, err := utils.GenerateJwt(strconv.Itoa(int(user.ID)), user.Role)
	if err != nil {
		log.Println("Failed to generate token:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to generate token",
		})
	}

	// Set the token as a cookie
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)

	// Redirect based on role
	if user.Role == "admin" || user.Role == "superadmin" {
		return c.Redirect("/dashboard", fiber.StatusFound)
	}

	return c.Redirect("/profile", fiber.StatusFound)

	// return c.JSON(fiber.Map{
	// 	"token": token,
	// })
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
