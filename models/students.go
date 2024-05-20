package models

import "gorm.io/gorm"

type student struct {
	*gorm.Model
	Batch        string `gorm:"not null" json:"batch"`
	Program      string `gorm:"not null" json:"program"`
	Symbol       string `gorm:"not null" json:"symbol"`
	Registration string `gorm:"not null" json:"registration"`
	Fullname     string `gorm:"not null" json:"fullname"`
}
