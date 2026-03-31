package repositories

import (
	"errors"

	"github/folkyyyy/preorder-api/internal/apperrors"
	"github/folkyyyy/preorder-api/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type OrderRepository interface {
	CreateOrder(order *models.Order, items []models.OrderItem) error
	GetOrdersByRoundID(roundID uint) ([]models.Order, error)
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db}
}

func (r *orderRepository) CreateOrder(order *models.Order, items []models.OrderItem) error {
	// 1. เริ่ม Transaction
	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// 2. สร้างบิล Order ลงตารางก่อน เพื่อให้ได้ order.ID มาใช้งาน
	if err := tx.Create(order).Error; err != nil {
		tx.Rollback()
		return err
	}

	var totalAmount float64

	// 3. วนลูปจัดการ OrderItem ทีละรายการ
	for i := range items {
		var preorderMenu models.PreorderMenu

		// 🌟 ไฮไลท์สำคัญ: ใช้ clause.Locking{Strength: "UPDATE"} เพื่อล็อก Row นี้ไว้!
		// แปลว่าถ้ามี Request อื่นพยายามจะอ่าน/แก้ เมนูตัวนี้ มันจะต้อง "รอ" จนกว่าเราจะทำเสร็จ
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ? AND preorder_round_id = ?", items[i].PreorderMenuID, order.PreorderRoundID).
			First(&preorderMenu).Error; err != nil {
			
			tx.Rollback()
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperrors.ErrMenuNotFound
			}
			return err
		}

		// 4. เช็คว่า "ของเหลือพอไหม?" (สั่งไปแล้ว + กำลังจะสั่ง > โควต้าทั้งหมด)
		if preorderMenu.OrderedCount+items[i].Quantity > preorderMenu.Quota {
			tx.Rollback() // พัง! ยกเลิกบิลนี้ทิ้งทันที
			return apperrors.ErrQuotaExceeded
		}

		// 5. ดึงข้อมูลราคาจากตาราง Menu หลัก (ป้องกันแอดมินส่งราคาหลอกมา)
		var menu models.Menu
		if err := tx.First(&menu, preorderMenu.MenuID).Error; err != nil {
			tx.Rollback()
			return errors.New("เกิดข้อผิดพลาดในการดึงข้อมูลราคาเมนู")
		}

		// 6. ประกอบร่าง OrderItem และคำนวณยอดรวม
		items[i].OrderID = order.ID
		items[i].PriceAtOrder = menu.Price // เอาราคา ณ ปัจจุบันมาใส่
		totalAmount += items[i].PriceAtOrder * float64(items[i].Quantity)

		// 7. อัปเดตยอดสั่งซื้อ (OrderedCount) กลับไปที่ตาราง PreorderMenu
		if err := tx.Model(&preorderMenu).
			UpdateColumn("ordered_count", preorderMenu.OrderedCount+items[i].Quantity).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// 8. บันทึก OrderItems ทั้งหมดลงฐานข้อมูลรวดเดียว
	if err := tx.Create(&items).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 9. อัปเดตยอดรวม (TotalAmount) กลับไปที่บิล Order หลัก
	if err := tx.Model(order).UpdateColumn("total_amount", totalAmount).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 10. ทุกอย่างผ่านฉลุย ยืนยันการบันทึก!
	return tx.Commit().Error
}


// สำหรับดึงบิลทั้งหมดในรอบนั้นๆ มาแสดงในหน้ารายการบิล (Order List) ของแอดมิน
func (r *orderRepository) GetOrdersByRoundID(roundID uint) ([]models.Order, error) {
	var orders []models.Order

	err := r.db.
		Preload("User"). // ดึงข้อมูล User (ถ้ามี, ถ้าเป็น NULL มันก็จะไม่พังครับ)
		Preload("OrderItems"). // 1. ดึงรายการอาหารในบิล
		Preload("OrderItems.PreorderMenu"). // 2. ดึงข้อมูลว่าสั่งจากโควต้าไหน
		Preload("OrderItems.PreorderMenu.Menu"). // 3.  ดึงชื่ออาหารและรูปภาพจากตารางเมนูหลักมาโชว์!
		Where("preorder_round_id = ?", roundID).
		Order("created_at DESC"). // เรียงบิลใหม่ล่าสุดไว้บนสุด
		Find(&orders).Error

	return orders, err
}