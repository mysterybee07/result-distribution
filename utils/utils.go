package utils

import (
	"errors"
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

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
		HTTPOnly: true,   // Prevents JavaScript from accessing the cookie
		Secure:   false,  // Set to true for HTTPS in production
		SameSite: "None", // Adjust based on your needs (e.g., "Strict" or "None")
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

func RandLetter(n int) string {
	b := make([]rune, n)
	rand.Seed(time.Now().UnixNano())
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func SanitizeFileName(fileName string) string {
	// Remove any characters that are not alphanumeric, dot, or underscore
	re := regexp.MustCompile(`[^a-zA-Z0-9._-]`)
	return strings.ToLower(re.ReplaceAllString(fileName, "_"))
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

// ConvertIDs converts batch ID and program ID from string to pointers of uint.
// Returns an error if any conversion fails.
func ConvertIDs(batchIDStr, programIDStr string) (batchID, programID *uint, err error) {
	if batchIDStr != "" {
		batchUint, err := strconv.ParseUint(batchIDStr, 10, 32)
		if err != nil {
			return nil, nil, err // Return the error if conversion fails
		}
		batchIDVal := uint(batchUint)
		batchID = &batchIDVal // Return a pointer to the converted uint
	}

	if programIDStr != "" {
		programUint, err := strconv.ParseUint(programIDStr, 10, 32)
		if err != nil {
			return nil, nil, err // Return the error if conversion fails
		}
		programIDVal := uint(programUint)
		programID = &programIDVal // Return a pointer to the converted uint
	}

	return batchID, programID, nil // Return both pointers and nil error
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
