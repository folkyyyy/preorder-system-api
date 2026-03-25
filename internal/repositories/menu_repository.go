package repositories

import (
	"github/folkyyyy/preorder-api/internal/models"
	"gorm.io/gorm"
)

// 1. สร้าง Interface กำหนดว่า Repository นี้ทำอะไรได้บ้าง
type MenuRepository interface {
	CreateMenu(menu *models.Menu)error
	GetAllMenus() ([]models.Menu, error)
}

type menuRepository struct {
	db *gorm.DB
}