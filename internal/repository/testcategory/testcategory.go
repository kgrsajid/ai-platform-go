package testcategory

import (
	"project-go/internal/models"

	"gorm.io/gorm"
)

type TestCategoryRepository struct {
	db *gorm.DB
}

func New(db *gorm.DB) *TestCategoryRepository {
	return &TestCategoryRepository{db: db}
}

func (r *TestCategoryRepository) CreateCategory(category *models.Category) (*models.Category, error) {
	if err := r.db.Create(&category).Error; err != nil {
		return nil, err
	}
	return category, nil
}

func (r *TestCategoryRepository) GetAllCategory() ([]models.Category, error) {
	var categories []models.Category
	err := r.db.Find(&categories).Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}
