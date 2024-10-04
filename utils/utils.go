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
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

var JwtSecret = []byte("Ajfdslfjlsdfjldslfj")

func GenerateJwt(userID uint, role string, c *fiber.Ctx) (string, error) {
	expirationTime := time.Now().Add(5 * time.Minute) // Set token expiration time
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
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    tokenString,
		Expires:  expirationTime,
		HTTPOnly: true, // Prevent JavaScript from accessing the cookie
	}

	c.Cookie(&cookie)
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
