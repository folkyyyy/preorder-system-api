package models

import "gorm.io/gorm"

type User struct {
	gorm.Model        // ใส่บรรทัดนี้ GORM จะสร้างฟิลด์ ID, CreatedAt, UpdatedAt, DeletedAt ให้เราอัตโนมัติ
	Email      string `gorm:"uniqueIndex;not null"`
	Username   string
	Password   string `gorm:"not null"`
	Role       string `gorm:"default:'user'"` // เอาไว้แยก 'admin' หรือ 'user'
}
