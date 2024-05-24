// utils/jwt.go
package utils

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtSecret = []byte("Ajfdslfjlsdfjldslfj")

func GenerateJwt(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userID,
		"exp":    time.Now().Add(time.Hour * 24).Unix(),
	})

	return token.SignedString(jwtSecret)
}

func ParseJwt(tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", err
	}

	userID, ok := claims["userID"].(string)
	if !ok {
		return "", err
	}

	return userID, nil
}
