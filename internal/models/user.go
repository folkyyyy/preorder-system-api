package models

import "gorm.io/gorm"

type User struct {
	gorm.Model        // ใส่บรรทัดนี้ GORM จะสร้างฟิลด์ ID, CreatedAt, UpdatedAt, DeletedAt ให้เราอัตโนมัติ
	Email      string `gorm:"uniqueIndex;not null" json:"email"` // ต้องไม่ซ้ำกัน และห้ามเป็นค่าว่าง
	Username   string `json:"username"`                          // ต้องไม่ซ้ำกัน และห้ามเป็นค่าว่าง
	Password   string `gorm:"not null" json:"password"`          // ห้ามเป็นค่าว่าง
	Role       string `gorm:"default:'user'" json:"role"`        // กำหนดค่าเริ่มต้นเป็น "user"
}
