package testhandler

import (
	"log/slog"
	"net/http"

	"project-go/internal/dto/request"
	res "project-go/internal/dto/response"
	"project-go/internal/lib/auth"
	"project-go/internal/lib/response"
	testservice "project-go/internal/service/test"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type testDetailResponse struct {
	response.Response
	Test res.TestDetailsResponse `json:"data"`
}

func Create(log *slog.Logger, svc *testservice.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.test.Create"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		userId, ok := auth.GetUserID(r)
		if !ok {
			log.Error("unauthorized: missing user id")
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, response.Error("unauthorized"))
			return
		}

		var req request.TestRequest
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("failed to decode request", slog.String("error", err.Error()))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("invalid request body"))
			return
		}
		req.AuthorId = userId

		if req.Title == "" || req.Difficulty == "" || len(req.Categories) == 0 {
			log.Error("missing required fields")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("title, difficulty, and categories are required"))
			return
		}

		createdTest, err := svc.TestCreate(req)
		if err != nil {
			log.Error("failed to create test", slog.String("error", err.Error()))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to create test"))
			return
		}

		questionsResp := mapQuestionsToResponse(createdTest.Questions)

		render.JSON(w, r, testDetailResponse{
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
