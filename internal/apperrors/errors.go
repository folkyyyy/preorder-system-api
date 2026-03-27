package apperrors

import "errors"

// รวม Error ของระบบ Preorder Round
var (
	ErrPastDate      = errors.New("วันที่จัดส่งต้องไม่ใช่วันที่ผ่านไปแล้ว")
	ErrRoundConflict = errors.New("มีรอบพรีออเดอร์เปิดรับในวันที่นี้แล้ว ไม่สามารถเปิดซ้ำได้")
	ErrInvalidQuota  = errors.New("โควต้าอาหารต้องมากกว่า 0")
	ErrInvalidMenu   = errors.New("กรุณาเพิ่มเมนูอย่างน้อยหนึ่งรายการในรอบพรีออเดอร์")
)

// Menu-related errors
var (
	ErrPriceNegative = errors.New("ราคาอาหารต้องไม่ติดลบ")
)

// auth
var (
	ErrInvalidCredentials = errors.New("ชื่อผู้ใช้หรือรหัสผ่านไม่ถูกต้อง")
)
