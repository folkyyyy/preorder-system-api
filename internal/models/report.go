package models


// 1. ตัวรับข้อมูลดิบจาก Database (ที่แยกบรรทัดตามโน๊ตมาแล้ว)
type RawKitchenSummary struct {
	MenuID        uint
	MenuName      string
	IsSpecial     bool
	Note          string
	TotalQuantity int
	TotalRevenue  float64
}

// 2. ตัวแพ็คเกจที่จะส่งเป็น JSON ให้ Frontend (ตัวนี้ยอดรวมแล้ว และโน๊ตเป็น Array)
type KitchenSummary struct {
	MenuID        uint     `json:"menuId"`
	MenuName      string   `json:"menuName"`
	IsSpecial     bool     `json:"isSpecial"`
	TotalQuantity int      `json:"totalQuantity"`
	TotalRevenue  float64  `json:"totalRevenue"`
	Notes         []string `json:"notes"` // 🌟 เปลี่ยนเป็น Array เพื่อเก็บลิสต์ข้อความ
}
