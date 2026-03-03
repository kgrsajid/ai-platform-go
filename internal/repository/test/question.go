package testrepo

import (
	"project-go/internal/models"

	"gorm.io/gorm"
)

type QuestionRepository struct {
	db *gorm.DB
}

func NewQuestionRepo(db *gorm.DB) *QuestionRepository {
	return &QuestionRepository{db: db}
}

func (r *QuestionRepository) CreateQuestion(question *models.TestQuestion) (*models.TestQuestion, error) {
	if err := r.db.Create(&question).Error; err != nil {
		return nil, err
	}
	return question, nil
}
