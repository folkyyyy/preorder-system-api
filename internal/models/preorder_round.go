package models

import (
	"gorm.io/gorm"
	"time"
)

type Status string

// 2. กำหนดค่า Constants ที่อนุญาตให้ใช้ได้
const (
	StatusOpen Status = "open"
	StatusClosed Status = "closed"
)

// รอบการรับพรีออเดอร์ (เช่น รอบวันเสาร์ที่ 28)
type PreorderRound struct {
	gorm.Model
	Title        string    `gorm:"not null" json:"title"`        // ชื่อรอบ เช่น "รอบส่งวันศุกร์"
	DeliveryDate time.Time `gorm:"not null" json:"deliveryDate"` // วันที่จัดส่ง/รับอาหาร
	Status       Status    `gorm:"default:'open'" json:"status"` // open, closed
	// 1 รอบ มีได้หลายเมนู
	PreorderMenus []PreorderMenu `gorm:"foreignKey:PreorderRoundID"`
}

// ตารางตรงกลาง: จัดการว่า "รอบนี้" มี "เมนูอะไรบ้าง" และรับ "อย่างละเท่าไหร่"
type PreorderMenu struct {
	gorm.Model
	PreorderRoundID uint `gorm:"not null" json:"preorderRoundId"`
	MenuID          uint `gorm:"not null" json:"menuId"`
	Menu            Menu `gorm:"foreignKey:MenuID"`             // เชื่อมไปดึงข้อมูลเมนูหลัก
	Quota           int  `gorm:"not null" json:"quota"`         // จำนวนที่รับพรีออเดอร์ (อย่างละเท่าไหร่)
	OrderedCount    int  `gorm:"default:0" json:"orderedCount"` // สั่งไปแล้วกี่ที่ (เอาไว้เช็คของหมด)
}
