package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User represents the structure of the registration form data
type User struct {
	*gorm.Model
	Batch        int    `gorm:"not null" json:"batch"`
	Symbol       string `gorm:"type:varchar(100);not null" json:"symbol"`
	Registration string `gorm:"type:varchar(100);not null" json:"registration"`
	Fullname     string `gorm:"type:varchar(100);not null" json:"fullname"`
	Email        string `gorm:"type:varchar(100);unique;not null" json:"email"`
	Password     string `gorm:"type:varchar(100);not null" json:"password"`
	Terms        bool   `gorm:"not null" json:"terms"`
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
