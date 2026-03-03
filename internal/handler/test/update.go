package testhandler

import (
	"log/slog"
	"net/http"
	"strconv"

	"project-go/internal/dto/request"
	res "project-go/internal/dto/response"
	"project-go/internal/lib/auth"
	"project-go/internal/lib/response"
	testservice "project-go/internal/service/test"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

func Update(log *slog.Logger, svc *testservice.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.test.Update"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req request.TestRequest
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("failed to decode request")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("invalid request body"))
			return
		}

		testIdStr := chi.URLParam(r, "testId")
		if testIdStr == "" {
			log.Error("missing test id")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("missing test id"))
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
			log.Error("unauthorized: missing user id")
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, response.Error("unauthorized"))
			return
		}
		req.AuthorId = userId

		if req.Title == "" || req.Difficulty == "" || len(req.Categories) == 0 {
			log.Error("missing required fields")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("title, difficulty, and categories are required"))
			return
		}

		updatedTest, err := svc.TestUpdate(req, uint(testId))
		if err != nil {
			log.Error("failed to update test", slog.String("error", err.Error()))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to update test"))
			return
		}

		questionsResp := mapQuestionsToResponse(updatedTest.Questions)

		render.JSON(w, r, testDetailResponse{
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
