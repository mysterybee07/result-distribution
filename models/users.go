package models

import (
	"gorm.io/gorm"
)

// User represents the structure of the registration form data
type User struct {
	gorm.Model
	BatchID            *uint    `gorm:"type:bigint;index" json:"batch_id,omitempty"`   // Nullable BatchID
	ProgramID          *uint    `gorm:"type:bigint;index" json:"program_id,omitempty"` // Nullable ProgramID
	SymbolNumber       string   `gorm:"type:varchar(100);not null" json:"symbol_number"`
	RegistrationNumber string   `gorm:"type:varchar(100);not null" json:"registration_number"`
	Email              string   `gorm:"type:varchar(100);unique;not null" json:"email"`
	Password           string   `gorm:"type:varchar(100);not null" json:"-"`
	Role               string   `gorm:"type:varchar(20);default:user" json:"role"`
	ImageURL           string   `gorm:"type:varchar(255)" json:"image_url,omitempty"`
	Batch              *Batch   `gorm:"foreignkey:BatchID;constraint:OnDelete:SET NULL;"`   // Nullable foreign key
	Program            *Program `gorm:"foreignkey:ProgramID;constraint:OnDelete:SET NULL;"` // Nullable foreign key
}

type UserInput struct {
	ProgramID          uint   `form:"program_id" json:"program_id"`
	BatchID            uint   `form:"batch_id" json:"batch_id"`
	SymbolNumber       string `form:"symbol_number" json:"symbol_number"`
	RegistrationNumber string `form:"registration_number" json:"registration_number"`
	Email              string `form:"email" json:"email"`
	Password           string `form:"password" json:"password"`
	Role               string `form:"role" json:"role"`
}
