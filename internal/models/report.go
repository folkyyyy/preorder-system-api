package models

type KitchenSummary struct {
	MenuID        uint    `json:"menu_id"`
	MenuName      string  `json:"menu_name"`
	TotalQuantity int     `json:"total_quantity"`
	TotalRevenue  float64 `json:"total_revenue"` // เผื่อแอดมินอยากดูว่าเมนูนี้ทำเงินไปเท่าไหร่
	IsSpecial     bool    `json:"isSpecial"`
}
