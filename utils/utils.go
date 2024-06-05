package utils

import (
	"errors"
	"math/rand"
	"regexp"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

var jwtSecret = []byte("Ajfdslfjlsdfjldslfj")

func GenerateJwt(userID string, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userID,
		"role":   role,
		"exp":    time.Now().Add(24 * time.Hour).Unix(),
	})
	return token.SignedString(jwtSecret)
}

func ParseJwt(tokenStr string) (string, string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
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

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandLetter(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func SanitizeFileName(fileName string) string {
	return strings.Map(func(r rune) rune {
		if r == '/' || r == '\\' {
			return -1
		}
		return r
	}, fileName)
}
