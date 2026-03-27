package services

import (
	"github/folkyyyy/preorder-api/internal/apperrors"
	"github/folkyyyy/preorder-api/internal/models"
	"github/folkyyyy/preorder-api/internal/repositories"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(user *models.User) error
	Login(emailOrUserName string, password string) (string, error)
}

type authService struct {
	repo repositories.UserRepository
}

func NewAuthService(repo repositories.UserRepository) AuthService {
	return &authService{repo}
}

// Register: นำรหัสผ่านไปเข้ารหัส (Hash) ก่อนบันทึกลง DB
func (s *authService) Register(user *models.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	return s.repo.CreateUser(user)
}

// Login: ตรวจสอบข้อมูลผู้ใช้และสร้าง JWT token
func (s *authService) Login(emailOrUserName string, password string) (string, error) {
	user, err := s.repo.GetUserByEmailOrUserName(emailOrUserName, emailOrUserName)
	if err != nil {
		return "", apperrors.ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", apperrors.ErrInvalidCredentials
	}


	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // Token หมดอายุใน 1 วัน
	})	

	// เซ็น Token ด้วย JWT_SECRET
	secret := os.Getenv("JWT_SECRET")
	t, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return t, nil
}