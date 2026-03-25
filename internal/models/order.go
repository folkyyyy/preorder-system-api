package models

import "gorm.io/gorm"

// ข้อมูลหลักของใบออเดอร์ลูกค้า
type Order struct {
	gorm.Model
	UserID           uint          `gorm:"not null"`
	Name			 string        `gorm:"not null"` // ชื่อผู้รับ
	PreorderRoundID  uint          `gorm:"not null"` // ออเดอร์นี้สั่งในรอบไหน
	PreorderRound    PreorderRound `gorm:"foreignKey:PreorderRoundID"`
	DeliveryLocation string
	TotalAmount      float64
	Status           string        `gorm:"default:'pending'"` // pending, paid, completed
	
	OrderItems []OrderItem `gorm:"foreignKey:OrderID"`
}

// รายการอาหารที่อยู่ในใบออเดอร์
type OrderItem struct {
	gorm.Model
	OrderID        uint         `gorm:"not null"`
	PreorderMenuID uint         `gorm:"not null"` // ชี้ไปที่เมนูของรอบนั้นๆ
	PreorderMenu   PreorderMenu `gorm:"foreignKey:PreorderMenuID"`
	Quantity       int          `gorm:"not null"`
	Note           string       // หมายเหตุ เช่น ไม่ใส่กระเทียม
	PriceAtOrder   float64      // เก็บราคา ณ ตอนที่สั่ง (เผื่ออนาคตเมนูหลักขึ้นราคา บิลเก่าจะได้ไม่เพี้ยน)
}