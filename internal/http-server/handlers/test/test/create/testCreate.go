package test_create

import (
	"log/slog"
	"net/http"
	req "project-go/internal/http-server/dto/request"
	res "project-go/internal/http-server/dto/response"
	testservice "project-go/internal/http-server/service/test"
	"project-go/internal/lib/api/response"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Response struct {
	response.Response
	Test res.TestDetailsResponse `json:"data"`
}

func New(log *slog.Logger, testService *testservice.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.test.test.New"
		log = log.With(
			slog.String("op", op),
			slog.String("req_id", middleware.GetReqID(r.Context())),
		)
		var req req.TestRequest
		err := render.DecodeJSON(r.Body, &req)

		if err != nil {
			log.Error("failed to decode req", slog.String("error", err.Error()))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("failed to decode req"))
			return
		}
		if req.Title == "" || req.Difficulty == "" || len(req.Categories) == 0 {
			log.Error("some field is empty")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("some field is empty"))
			return
		}
		createdTest, err := testService.TestCreate(req)
		if err != nil {
			log.Error("failed to create test", slog.String("error", err.Error()))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to create test"))
			return
		}
		var questionsResp []res.QuestionResponse
		for _, q := range createdTest.Questions {
			var optionsResp []res.OptionResponse
			for _, o := range q.Options {
				optionsResp = append(optionsResp, res.OptionResponse{
					OptionText: o.OptionText,
					IsCorrect:  o.IsCorrect,
				})
			}

			// 2. Маппим сам вопрос
			questionsResp = append(questionsResp, res.QuestionResponse{
				Question: q.Question,
				Options:  optionsResp,
			})
		}

		render.JSON(w, r, Response{
			Response: response.OK(),
			Test: res.TestDetailsResponse{
				ID:          createdTest.ID,
				Title:       createdTest.Title,
				Description: createdTest.Description,
				CreatedAt:   createdTest.CreatedAt,
				Questions:   questionsResp,
				Difficulty:  string(createdTest.Difficulty),
				Tags:        createdTest.Tags,
				Categories:  createdTest.Categories,
			},
		})

	}
}
