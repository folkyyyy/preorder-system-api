package models

import "gorm.io/gorm"

type Role string

// 2. กำหนดค่า Constants ที่อนุญาตให้ใช้ได้
const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
	// ถ้าอนาคตมี Role อื่นก็มาเพิ่มตรงนี้ได้ เช่น RoleManager Role = "manager"
)

type User struct {
	gorm.Model         // ใส่บรรทัดนี้ GORM จะสร้างฟิลด์ ID, CreatedAt, UpdatedAt, DeletedAt ให้เราอัตโนมัติ
	Email      string  `gorm:"uniqueIndex;not null" json:"email"` // ต้องไม่ซ้ำกัน และห้ามเป็นค่าว่าง
	Username   string  `json:"username"`
	Password   string  `gorm:"not null" json:"password"` // ห้ามเป็นค่าว่าง
	Role       Role    `gorm:"type:varchar(20); default:'user'" json:"role"`
	Orders     []Order `gorm:"foreignKey:UserID"`
}
