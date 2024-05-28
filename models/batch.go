package models

import "gorm.io/gorm"

type Batch struct {
	gorm.Model
	Year uint `gorm:"not null" json:"year" form:"year"`
	// Programs []Program `gorm:"foreignKey:BatchID"`
}
