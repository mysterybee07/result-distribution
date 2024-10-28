package models

import (
	"gorm.io/gorm"
)

// College model
type College struct {
	gorm.Model
	CollegeCode string  `json:"college_code" gorm:"primaryKey;type:varchar(255);unique;not null"`
	CollegeName string  `json:"college_name" gorm:"not null"`
	Address     string  `json:"address" gorm:"not null"`
	Latitude    float64 `json:"latitude" gorm:"not null"`
	Longitude   float64 `json:"longitude" gorm:"not null"`
	IsCenter    bool    `json:"is_center" gorm:"default:false"` // Indicates if the college is registered as a center
	// Program     Program `gorm:"foreignKey:ProgramID"`           // Association with Program
	// Batch       Batch   `gorm:"foreignKey:BatchID"`             // Association with Batch
}

type CapacityAndCount struct {
	gorm.Model
	CollegeID     uint    `gorm:"not null" json:"college_id"`
	BatchID       uint    `gorm:"not null" json:"batch_id"`
	ProgramID     uint    `gorm:"not null" json:"program_id"`
	StudentsCount int     `gorm:"not null;default:0" json:"students_count"`
	Capacity      int     `json:"capacity" gorm:"default:0"` // Capacity if registered as a center
	College       College `gorm:"foreignKey:CollegeID"`
	Batch         Batch   `gorm:"foreignKey:BatchID"`
	Program       Program `gorm:"foreignKey:ProgramID"`
}

// Center model
// type Center struct {
// 	gorm.Model
// 	CenterCollegeCode string  `json:"center_college_code" gorm:"type:varchar(255);not null"`
// 	Name              string  `json:"name" gorm:"not null"`
// 	Address           string  `json:"address" gorm:"not null"`
// 	Latitude          float64 `json:"latitude" gorm:"not null"`
// 	Longitude         float64 `json:"longitude" gorm:"not null"`
// 	Capacity          int     `json:"capacity" gorm:"not null"`
// 	College           College `gorm:"foreignKey:CenterCollegeCode;references:CollegeCode"`
// }

// // Preference model
// type Preference struct {
// 	gorm.Model
// 	CollegeCode       string `json:"college_code" gorm:"not null"`
// 	CenterCollegeCode string `json:"center_college_code" gorm:"not null"`
// 	Preference        int    `json:"preference" gorm:"not null"`
// 	Reason            string `json:"reason"`
// }

// type AllocatedCenter struct {
// 	gorm.Model
// 	CollegeCode       string `json:"college_code" gorm:"not null"`
// 	CenterCollegeCode string `json:"center_college_code" gorm:"not null"`
// 	Count             int    `json:"count" gorm:"not null"`
// }
