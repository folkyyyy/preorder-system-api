package handlers

import (
	"github/folkyyyy/preorder-api/internal/models"
	"github/folkyyyy/preorder-api/internal/services"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	service services.AuthService
}

func NewAuthHandler(service services.AuthService) *AuthHandler {
	return &AuthHandler{service}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ส่งข้อมูลมาไม่ถูกต้อง"})
	}

	if err := h.service.Register(&user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "ไม่สามารถลงทะเบียนผู้ใช้ได้"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "ลงทะเบียนผู้ใช้สำเร็จ"})
}

// โครงสร้างรับข้อมูล Login
type LoginInput struct {
	EmailOrUserName string `json:"emailOrUserName"`
	Password string `json:"password"`
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var input LoginInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ส่งข้อมูลมาไม่ถูกต้อง"})
	}

	token, err := h.service.Login(input.EmailOrUserName, input.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message": "เข้าสู่ระบบสำเร็จ",
		"token":   token,
	})
}