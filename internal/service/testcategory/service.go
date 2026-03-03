package testcategoryservice

import (
	"project-go/internal/dto/request"
	"project-go/internal/models"
)

type TestCategoryRepository interface {
	CreateCategory(category *models.Category) (*models.Category, error)
	GetAllCategory() ([]models.Category, error)
}

type Service struct {
	categoryRepo TestCategoryRepository
}

func New(categoryRepo TestCategoryRepository) *Service {
	return &Service{categoryRepo: categoryRepo}
}

func (s *Service) CreateCategory(req request.TestCategoryReq) (*models.Category, error) {
	category := &models.Category{Name: req.Name}
	return s.categoryRepo.CreateCategory(category)
}

func (s *Service) GetAllCategories() ([]models.Category, error) {
	return s.categoryRepo.GetAllCategory()
}
