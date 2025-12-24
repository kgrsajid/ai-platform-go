package res

import (
	"project-go/internal/models"
	"time"
)

type TestDetailsResponse struct {
	ID          uint               `json:"id"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	Categories  []models.Category  `json:"categories"`
	Difficulty  string             `json:"difficulty"`
	Tags        []string           `json:"tags"`
	Questions   []QuestionResponse `json:"questions"`
	CreatedAt   time.Time          `json:"createdAt"`
}

func ToTestDetailsResponse(t *models.Test) *TestDetailsResponse {
	questions := make([]QuestionResponse, 0, len(t.Questions))
	for _, val := range t.Questions {
		options := make([]OptionResponse, 0, len(val.Options))
		for _, op := range val.Options {
			options = append(options, OptionResponse{
				ID:         op.ID,
				OptionText: op.OptionText,
				IsCorrect:  op.IsCorrect,
			})
		}
		questions = append(questions, QuestionResponse{
			Question: val.Question,
			Options:  options,
		})
	}
	return &TestDetailsResponse{
		ID:          t.ID,
		Title:       t.Title,
		Description: t.Description,
		Categories:  t.Categories,
		Difficulty:  string(t.Difficulty),
		Tags:        t.Tags,
		Questions:   questions,
		CreatedAt:   t.CreatedAt,
	}
}

type TestResponse struct {
	ID                uint           `json:"id"`
	Title             string         `json:"title"`
	Description       string         `json:"description"`
	Categories        []TestCategory `json:"categories"`
	Difficulty        string         `json:"difficulty"`
	NumberOfQuestions uint           `json:"numberOfQuestion"`
	Tags              []string       `json:"tags"`
	CreatedAt         time.Time      `json:"createdAt"`
}

type TestWithQuestionsCount struct {
	models.Test
	NumberOfQuestions int `gorm:"column:number_of_questions"`
}

func ToTestResponse(t TestWithQuestionsCount) TestResponse {
	return TestResponse{
		ID:                t.ID,
		Title:             t.Title,
		Difficulty:        string(t.Difficulty),
		Tags:              t.Tags,
		NumberOfQuestions: uint(t.NumberOfQuestions),
		Description:       t.Description,
		Categories:        toTestCategories(t.Categories),
	}
}

func toTestCategories(t []models.Category) []TestCategory {
	categroies := make([]TestCategory, 0, len(t))
	for _, category := range t {
		categroies = append(categroies, TestCategory{
			ID:   category.ID,
			Name: category.Name,
		})
	}
	return categroies
}

type QuestionResponse struct {
	Question string           `json:"question"`
	Options  []OptionResponse `json:"options"`
}

type OptionResponse struct {
	ID         uint   `json:"id"`
	OptionText string `json:"optionText"`
	IsCorrect  bool   `json:"isCorrect"`
}
