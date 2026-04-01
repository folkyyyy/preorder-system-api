package models

import "gorm.io/gorm"

type Menu struct {
	gorm.Model
	Name        string  `gorm:"not null" json:"name"`
	Description string  `json:"description"`
	Price       float64 `gorm:"not null" json:"price"`
	ImageURL    string  `json:"imageURL"` // เผื่อใส่รูปภาพในอนาคต

	IsSpecialAllowed bool    `gorm:"default:false" json:"isSpecialAllowed"`
	SpecialPrice     float64 `json:"specialPrice"`
}
