package testcategoryservice

import (
	req "project-go/internal/http-server/dto/request"
	"project-go/internal/models"
)

type TestCategoryRepo interface {
	CreateCategory(category *models.Category) (*models.Category, error)
	GetAllCategory() ([]models.Category, error)
}

type Service struct {
	TestCategoryRepo TestCategoryRepo
}

func New(TestCategoryRepo TestCategoryRepo) *Service {
	return &Service{
		TestCategoryRepo: TestCategoryRepo,
	}
}

func (s *Service) CreateCategory(category req.TestCategoryReq) (*models.Category, error) {
	testcategory := &models.Category{
		Name: category.Name,
	}
	res, err := s.TestCategoryRepo.CreateCategory(testcategory)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *Service) GetAllCategories() ([]models.Category, error) {
	categories, err := s.TestCategoryRepo.GetAllCategory()
	if err != nil {
		return nil, err
	}
	return categories, nil
}
