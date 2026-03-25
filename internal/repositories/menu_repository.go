package repositories

import (
	"github/folkyyyy/preorder-api/internal/models" 
	"gorm.io/gorm"
)

// 1. สร้าง Interface กำหนดว่า Repository นี้ทำอะไรได้บ้าง
type MenuRepository interface {
	CreateMenu(menu *models.Menu) error
	GetAllMenus() ([]models.Menu, error)
}

type menuRepository struct {
	db *gorm.DB
}

// 2. ฟังก์ชันสำหรับสร้าง Repository Instance
func NewMenuRepository(db *gorm.DB) MenuRepository {
	return &menuRepository{db}
}

// 3. เขียนฟังก์ชันสร้างเมนูใหม่ (Insert ลง DB)
func (r *menuRepository) CreateMenu(menu *models.Menu) error {
	return r.db.Create(menu).Error
}

// 4. เขียนฟังก์ชันดึงเมนูทั้งหมด (Select จาก DB)
func (r *menuRepository) GetAllMenus() ([]models.Menu, error) {
	var menus []models.Menu
	err := r.db.Find(&menus).Error
	return menus, err
}