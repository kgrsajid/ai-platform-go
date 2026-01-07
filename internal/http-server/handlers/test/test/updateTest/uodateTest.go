package updatetest

import (
	"log/slog"
	"net/http"
	req "project-go/internal/http-server/dto/request"
	res "project-go/internal/http-server/dto/response"
	testservice "project-go/internal/http-server/service/test"
	"project-go/internal/lib/api/response"
	"project-go/internal/lib/auth"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Response struct {
	response.Response
	Test res.TestDetailsResponse `json:"data"`
}

func New(log *slog.Logger, testService *testservice.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.test.updateTest.New"
		log = log.With(
			slog.String("op", op),
			slog.String("Req_id", middleware.GetReqID(r.Context())),
		)
		var req req.TestRequest
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("failed to decode request"))
			return
		}
		testIdStr := chi.URLParam(r, "testId")
		if testIdStr == "" {
			log.Error("test id is empty")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("test id is empty"))
			return
		}
		testId, err := strconv.ParseUint(testIdStr, 10, 64)
		if err != nil {
			log.Error("invalid test id")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("invalid test id"))
			return
		}
		userId, ok := auth.GetUserID(r)
		if !ok {
			log.Error("user id is null")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("user id is null"))
			return
		}
		req.AuthorId = userId

		if req.Title == "" || req.Difficulty == "" || len(req.Categories) == 0 {
			log.Error("some field is empty")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("some field is empty"))
			return
		}

		updatedTest, err := testService.TestUpdate(req, uint(testId))
		if err != nil {
			log.Error("failed to update test")
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to update test"))
			return
		}
		var questionsResp []res.QuestionResponse
		for _, q := range updatedTest.Questions {
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
				ID:          updatedTest.ID,
				Title:       updatedTest.Title,
				Description: updatedTest.Description,
				CreatedAt:   updatedTest.CreatedAt,
				Questions:   questionsResp,
				Difficulty:  string(updatedTest.Difficulty),
				Tags:        updatedTest.Tags,
				Categories:  updatedTest.Categories,
			},
		})
	}
}
