package models

import (
	"time"
	"gorm.io/gorm"
)

// รอบการรับพรีออเดอร์ (เช่น รอบวันเสาร์ที่ 28)
type PreorderRound struct {
	gorm.Model
	Title        string    `gorm:"not null"` // ชื่อรอบ เช่น "รอบส่งวันศุกร์"
	DeliveryDate time.Time `gorm:"not null"` // วันที่จัดส่ง/รับอาหาร
	Status       string    `gorm:"default:'open'"` // open, closed, completed
	
	// 1 รอบ มีได้หลายเมนู
	PreorderMenus []PreorderMenu `gorm:"foreignKey:PreorderRoundID"`
}

// ตารางตรงกลาง: จัดการว่า "รอบนี้" มี "เมนูอะไรบ้าง" และรับ "อย่างละเท่าไหร่"
type PreorderMenu struct {
	gorm.Model
	PreorderRoundID uint    `gorm:"not null"`
	MenuID          uint    `gorm:"not null"`
	Menu            Menu    `gorm:"foreignKey:MenuID"` // เชื่อมไปดึงข้อมูลเมนูหลัก
	Quota           int     `gorm:"not null"` // จำนวนที่รับพรีออเดอร์ (อย่างละเท่าไหร่)
	OrderedCount    int     `gorm:"default:0"` // สั่งไปแล้วกี่ที่ (เอาไว้เช็คของหมด)
}