package testhandler

import (
	res "project-go/internal/dto/response"
	"project-go/internal/models"
)

func mapQuestionsToResponse(questions []models.TestQuestion) []res.QuestionResponse {
	questionsResp := make([]res.QuestionResponse, 0, len(questions))
	for _, q := range questions {
		optionsResp := make([]res.OptionResponse, 0, len(q.Options))
		for _, o := range q.Options {
			optionsResp = append(optionsResp, res.OptionResponse{
				OptionText: o.OptionText,
				IsCorrect:  o.IsCorrect,
			})
		}
		questionsResp = append(questionsResp, res.QuestionResponse{
			Question: q.Question,
			Options:  optionsResp,
		})
	}
	return questionsResp
}
