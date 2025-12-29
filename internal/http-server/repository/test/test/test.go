package test

import (
	req "project-go/internal/http-server/dto/request"
	res "project-go/internal/http-server/dto/response"
	"project-go/internal/models"
	"time"

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
	if err := r.db.Where("id = ?", testId).Preload("Questions").Preload("Questions.Options").Preload("Categories").Find(&test).Error; err != nil {
		return nil, err
	}
	return test, nil
}

func (r *TestRepository) AddTestResult(testReq req.TestResultReq) (*models.TestResult, error) {
	var lastAttempt models.TestResult
	r.db.
		Where("student_id = ? AND test_id = ?", testReq.UserId, testReq.TestId).
		Order("attempt DESC").
		First(&lastAttempt)

	newAttempt := 1
	if lastAttempt.ID != 0 {
		newAttempt = lastAttempt.Attempt + 1
	}
	test := &models.TestResult{
		StudentID:   testReq.UserId,
		TestID:      testReq.TestId,
		Score:       testReq.Score,
		MaxScore:    testReq.MaxScore,
		Percentage:  (float64(testReq.Score) / float64(testReq.MaxScore)) * 100,
		StartedAt:   testReq.StartedAt,
		FinishedAt:  testReq.FinishedAt,
		Attempt:     newAttempt,
		DurationSec: testReq.DurationSec,
	}
	if err := r.db.Create(&test).Error; err != nil {
		return nil, err
	}
	return test, nil
}

func (r *TestRepository) GetAllUserTestResults(filter req.GetALlTestResultsFilter) ([]models.TestResult, error) {
	var results []models.TestResult

	query := r.db.Model(&models.TestResult{}).
		Where("student_id = ?", filter.UserID).
		Where("test_id = ?", filter.TestID)

	// фильтр по дате
	switch filter.Duration {
	case "day":
		from := time.Now().AddDate(0, 0, -1)
		query = query.Where("created_at >= ?", from)
	case "week":
		from := time.Now().AddDate(0, 0, -7)
		query = query.Where("created_at >= ?", from)
	case "month":
		from := time.Now().AddDate(0, -1, 0)
		query = query.Where("created_at >= ?", from)
	}

	err := query.Find(&results).Error
	if err != nil {
		return nil, err
	}

	return results, nil
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
	if filter.IsPrivate != nil && *filter.IsPrivate {
		query = query.Where("is_private = true AND author_id = ?", filter.UserId)
	} else {
		query = query.Where("is_private = false")
	}

	if filter.Search != nil && *filter.Search != "" {
		like := "%" + *filter.Search + "%"
		query = query.Where("(tests.title ILIKE ? OR tags @> ARRAY[?])", like, *filter.Search)
	}

	if len(filter.Categories) > 0 {
		query = query.Joins("JOIN test_categories tc ON tc.test_id = tests.id").
			Where("tc.category_id IN ?", filter.Categories)
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
		tests[i] = res.TestWithQuestionsCount{
			Test:              t,
			NumberOfQuestions: len(t.Questions),
		}
	}

	return tests, total, nil
}
