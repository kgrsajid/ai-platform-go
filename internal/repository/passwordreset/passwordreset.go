package passwordreset

import (
	"project-go/internal/models"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(pr *models.PasswordReset) error {
	return r.db.Create(pr).Error
}

func (r *Repository) FindByEmailAndCode(email, code string) (*models.PasswordReset, error) {
	var pr models.PasswordReset
	err := r.db.Where("email = ? AND code = ?", email, code).First(&pr).Error
	return &pr, err
}

func (r *Repository) FindByToken(token string) (*models.PasswordReset, error) {
	var pr models.PasswordReset
	err := r.db.Where("token = ?", token).First(&pr).Error
	return &pr, err
}

func (r *Repository) UpdateToken(id uint, token string) error {
	return r.db.Model(&models.PasswordReset{}).Where("id = ?", id).Update("token", token).Error
}

func (r *Repository) DeleteByEmail(email string) error {
	return r.db.Where("email = ?", email).Delete(&models.PasswordReset{}).Error
}
