package models

import "gorm.io/gorm"

type Menu struct {
	gorm.Model
	Name        string  `gorm:"not null"`
	Description string
	Price       float64 `gorm:"not null"`
	ImageURL    string  // เผื่อใส่รูปภาพในอนาคต
}