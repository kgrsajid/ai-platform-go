package test

import (
	"fmt"
	req "project-go/internal/http-server/dto/request"
	res "project-go/internal/http-server/dto/response"
	"project-go/internal/models"

	"gorm.io/gorm"
)

type TestRepository struct {
	db *gorm.DB
}

func NewTestRepo(db *gorm.DB) *TestRepository {
	return &TestRepository{db: db}
}

func (r *TestRepository) CreateTest(test *models.Test) (*models.Test, error) {
	var categories []models.Category
	for _, c := range test.Categories {
		categories = append(categories, models.Category{ID: c.ID})
	}
	// Загружаем существующие категории из базы
	if len(categories) > 0 {
		if err := r.db.Where("id IN ?", extractIDs(categories)).Find(&categories).Error; err != nil {
			return nil, err
		}
	}
	test.Categories = categories

	if err := r.db.Create(&test).Error; err != nil {
		return nil, err
	}
	return test, nil
}

func extractIDs(cats []models.Category) []uint {
	ids := make([]uint, len(cats))
	for i, c := range cats {
		ids[i] = c.ID
	}
	return ids
}

func (r *TestRepository) GetTestById(testId uint64) (*models.Test, error) {
	var test *models.Test
	if err := r.db.Where("id = ?", testId).Find(&test).Error; err != nil {
		return nil, err
	}
	return test, nil
}

func (r *TestRepository) GetAllTest(
	filter req.TestFilter,
) ([]res.TestWithQuestionsCount, int64, error) {

	var testModels []models.Test
	query := r.db.Model(&models.Test{}).
		Preload("Categories").
		Preload("Questions") // чтобы len(t.Questions) работал

	// Фильтры
	if filter.Difficulty != nil && *filter.Difficulty != "" {
		query = query.Where("difficulty = ?", *filter.Difficulty)
	}

	if filter.Search != nil && *filter.Search != "" {
		like := "%" + *filter.Search + "%"
		query = query.Where("(tests.title ILIKE ? OR tags @> ARRAY[?])", like, *filter.Search)
	}

	if filter.Category != nil && *filter.Category != "" {
		query = query.
			Joins("JOIN test_categories tc ON tc.test_id = tests.id").
			Joins("JOIN categories c ON c.id = tc.category_id").
			Where("c.name = ?", *filter.Category)
	}

	// COUNT тестов
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Основной запрос с пагинацией
	if err := query.Order("tests.created_at DESC").
		Limit(filter.Limit).
		Offset(filter.Offset).
		Find(&testModels).Error; err != nil {
		return nil, 0, err
	}

	// Маппим в DTO
	tests := make([]res.TestWithQuestionsCount, len(testModels))
	for i, t := range testModels {
		fmt.Println(t.Categories, t.Description)
		tests[i] = res.TestWithQuestionsCount{
			Test:              t,
			NumberOfQuestions: len(t.Questions),
		}
	}

	return tests, total, nil
}
