package user

import (
	"project-go/internal/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(u *models.User) (*models.User, error) {
	if err := r.db.Create(&u).Error; err != nil {
		return nil, err
	}
	return u, nil
}

func (r *UserRepository) FindUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	return &user, err
}

func (r *UserRepository) UpdatePassword(id uint, hash string) error {
	return r.db.Model(&models.User{}).Where("id = ?", id).Update("password", hash).Error
}
