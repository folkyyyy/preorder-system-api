package middlewares

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// Protected เป็นด่านตรวจว่ามี Token ส่งมาไหม และ Token ถูกต้องไหม
func Protected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 1. ดึง Token จาก Header (Authorization: Bearer <token>)
		authHeader := c.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized: Missing or invalid token format"})
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		secret := os.Getenv("JWT_SECRET")

		// 2. ตรวจสอบและถอดรหัส Token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized: Invalid or expired token"})
		}

		// 3. เอาข้อมูลใน Token (เช่น user_id, role) ไปฝากไว้ใน Context เผื่อให้ API ปลายทางเอาไปใช้ต่อ
		claims := token.Claims.(jwt.MapClaims)
		c.Locals("user_id", claims["user_id"])
		c.Locals("role", claims["role"])

		// 4. ให้ผ่านไปทำฟังก์ชันถัดไปได้
		return c.Next()
	}
}

func AdminOnly() fiber.Handler {
	return func(c *fiber.Ctx) error {

		role, ok := c.Locals("role").(string)

		if !ok || role != "admin" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "คุณไม่มีสิทธิ์เข้าถึงทรัพยากรนี้",
			})
		}
		return c.Next()
	}

}
