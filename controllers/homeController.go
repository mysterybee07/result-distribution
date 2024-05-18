package controllers

import (
	"log"
	"regexp"
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

func validateEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	regex := regexp.MustCompile(pattern)
	return regex.MatchString(email)
}
func StoreRegister(c *fiber.Ctx) error {
	var data models.User

	// Parse the form data into the User struct
	if err := c.BodyParser(&data); err != nil {
		log.Println("Unable to parse form data:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid form data",
		})
	}

	// Validate fields
	if data.Batch == 0 ||
		data.Symbol == "" ||
		data.Registration == "" ||
		data.Fullname == "" ||
		data.Email == "" ||
		data.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "All fields are required",
		})
	}

	// Check if the user password is less than 8 characters
	if len(data.Password) < 8 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Password must be at least 8 characters long",
		})
	}

	// Validate email
	if !validateEmail(data.Email) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid email",
		})
	}

	// Check if email already exists
	var existingUser models.User
	if err := initializers.DB.Where("email = ?", data.Email).First(&existingUser).Error; err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Email is already taken",
		})
	}

	// Hash password
	hashedPassword, err := models.HashPassword(data.Password)
	if err != nil {
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
	return c.Redirect("/login")

	// return c.Status(fiber.StatusCreated).JSON(fiber.Map{
	// 	"user":    data,
	// 	"message": "Account created successfully",
	// })
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

	var data LoginData

	// Parse the form data into the LoginData struct
	if err := c.BodyParser(&data); err != nil {
		log.Println("Unable to parse form data:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid form data",
		})
	}

	log.Printf("Parsed login data: %+v\n", data)

	// Check if identifier (symbol number or email) and password are provided
	if data.Identifier == "" || data.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Identifier and password are required",
		})
	}

	// Attempt to find the user by email or symbol number
	var user models.User
	if err := initializers.DB.Where("email = ? OR symbol= ?", data.Identifier, data.Identifier).First(&user).Error; err != nil {
		log.Printf("User not found with identifier: %s\n", data.Identifier)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid identifier or password",
		})
	}

	log.Printf("Found user: %+v\n", user)

	// Verify the provided password against the stored hashed password
	if !models.CheckPasswordHash(data.Password, user.Password) {
		log.Println("Invalid password")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid identifier or password",
		})
	}

	// Generate JWT token
	token, err := utils.GenerateJwt(strconv.Itoa(int(user.ID)))
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

	// return c.Status(fiber.StatusOK).JSON(fiber.Map{
	// 	"message": "Successfully logged in",
	// 	"user":    user,
	// })
	return c.Redirect("/profile")
}
