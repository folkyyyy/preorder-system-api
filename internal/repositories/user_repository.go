package repositories

import(
	"github/folkyyyy/preorder-api/internal/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *models.User) error
	GetUserByEmailOrUserName(email string, userName string) (*models.User, error)
}

type userRepository struct{
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) CreateUser(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) GetUserByEmailOrUserName(email string, userName string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ? OR username = ?", email, userName).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}