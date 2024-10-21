package models

import "gorm.io/gorm"

// College model
type College struct {
	gorm.Model
	CollegeCode   string       `json:"college_code" gorm:"primaryKey;type:varchar(255);unique;not null"`
	Name          string       `json:"name" gorm:"not null"`
	Address       string       `json:"address" gorm:"not null"`
	Latitude      float64      `json:"latitude" gorm:"not null"`
	Longitude     float64      `json:"longitude" gorm:"not null"`
	StudentsCount int          `json:"students_count" gorm:"not null"`
	Centers       []Center     `gorm:"foreignKey:CenterCollegeCode;references:CollegeCode"`
	Preferences   []Preference `gorm:"foreignKey:CollegeCode;references:CollegeCode"`
}

// Center model
type Center struct {
	gorm.Model
	CenterCollegeCode string  `json:"center_college_code" gorm:"type:varchar(255);not null"`
	Name              string  `json:"name" gorm:"not null"`
	Address           string  `json:"address" gorm:"not null"`
	Latitude          float64 `json:"latitude" gorm:"not null"`
	Longitude         float64 `json:"longitude" gorm:"not null"`
	Capacity          int     `json:"capacity" gorm:"not null"`
	College           College `gorm:"foreignKey:CenterCollegeCode;references:CollegeCode"`
}

// Preference model
type Preference struct {
	gorm.Model
	CollegeCode       string `json:"college_code" gorm:"not null"`
	CenterCollegeCode string `json:"center_college_code" gorm:"not null"`
	Preference        int    `json:"preference" gorm:"not null"`
	Reason            string `json:"reason"`
}

type AllocatedCenter struct {
	gorm.Model
	CollegeCode       string `json:"college_code" gorm:"not null"`
	CenterCollegeCode string `json:"center_college_code" gorm:"not null"`
	Count             int    `json:"count" gorm:"not null"`
}
