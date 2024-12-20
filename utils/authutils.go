package utils

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

var JwtSecret = []byte("Ajfdslfjlsdfjldslfj")

func GenerateJwt(userID uint, role string, c *fiber.Ctx) (string, error) {
	expirationTime := time.Now().Add(time.Hour * 24) // Set token expiration time
	claims := jwt.MapClaims{
		"userID": strconv.Itoa(int(userID)),
		"role":   role,
		"exp":    expirationTime.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(JwtSecret) // Use the correct signing key (HS256)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	// Set JWT token as a cookie
	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    tokenString,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
		Secure:   false,
		SameSite: "None",
	})

	return tokenString, nil
}

func ParseJwt(tokenStr string) (string, string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return JwtSecret, nil
	})

	if err != nil || !token.Valid {
		return "", "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", "", err
	}

	userID, ok := claims["userID"].(string)
	if !ok {
		return "", "", errors.New("invalid token claims: userID not found")
	}

	role, ok := claims["role"].(string)
	if !ok {
		return "", "", errors.New("invalid token claims: role not found")
	}

	return userID, role, nil
}

func ValidateEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	regex := regexp.MustCompile(pattern)
	return regex.MatchString(email)
}

// HashPassword hashes a plain text password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash compares a plain text password with a hashed password
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
