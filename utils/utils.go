package utils

import (
	"time"

	"github.com/golang-jwt/jwt"
)

var jwtKey = []byte("jasldfjlsahdflshdfl")

// GenerateJwt generates a new JWT token
func GenerateJwt(userId string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userId,
		"exp":    time.Now().Add(24 * time.Hour).Unix(),
	})
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func Parsejwt(cookie string) (string, error) {
	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{},
		func(t *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		})
	if err != nil || token.Valid {
		return "", err
	}
	claims := token.Claims.(*jwt.StandardClaims)
	return claims.Issuer, nil
}
