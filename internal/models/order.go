package models

import "gorm.io/gorm"

// ข้อมูลหลักของใบออเดอร์ลูกค้า
type Order struct {
	gorm.Model
	UserID           uint          `gorm:"not null" json:"userId"`          // ใครเป็นคนสั่ง
	Name             string        `gorm:"not null" json:"name"`            // ชื่อผู้รับ
	PreorderRoundID  uint          `gorm:"not null" json:"preorderRoundId"` // ออเดอร์นี้สั่งในรอบไหน
	PreorderRound    PreorderRound `gorm:"foreignKey:PreorderRoundID"`
	DeliveryLocation string        `json:"deliveryLocation"`
	TotalAmount      float64       `json:"totalAmount"`
	Status           string        `gorm:"default:'pending'" json:"status"` // pending, paid, completed

	OrderItems []OrderItem `gorm:"foreignKey:OrderID"`
}

// รายการอาหารที่อยู่ในใบออเดอร์
type OrderItem struct {
	gorm.Model
	OrderID        uint         `gorm:"not null" json:"orderId"`        // ชี้ไปที่ใบออเดอร์หลัก
	PreorderMenuID uint         `gorm:"not null" json:"preorderMenuId"` // ชี้ไปที่เมนูของรอบนั้นๆ
	PreorderMenu   PreorderMenu `gorm:"foreignKey:PreorderMenuID"`
	Quantity       int          `gorm:"not null" json:"quantity"`
	Note           string       `json:"note"`         // หมายเหตุ เช่น ไม่ใส่กระเทียม
	PriceAtOrder   float64      `json:"priceAtOrder"` // เก็บราคา ณ ตอนที่สั่ง (เผื่ออนาคตเมนูหลักขึ้นราคา บิลเก่าจะได้ไม่เพี้ยน)
}
