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
	"github.com/go-chi/render"
)

type testResultsResponse struct {
	response.Response
	Data []res.TestResultResponse `json:"data"`
}

func GetAllUserResults(log *slog.Logger, svc *testservice.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId, ok := auth.GetUserID(r)
		if !ok {
			log.Error("unauthorized: missing user id")
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, response.Error("unauthorized"))
			return
		}

		testIdStr := chi.URLParam(r, "testId")
		testId, err := strconv.ParseUint(testIdStr, 10, 64)
		if err != nil {
			http.Error(w, "invalid testId", http.StatusBadRequest)
			return
		}

		duration := r.URL.Query().Get("duration")
		if duration == "" {
			duration = "month"
		}

		filter := request.GetAllTestResultsFilter{
			UserID:   uint(userId),
			TestID:   uint(testId),
			Duration: duration,
		}

		tests, err := svc.GetAllUserTestResults(filter)
		if err != nil {
			log.Error("failed to get test results", slog.String("error", err.Error()))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to get test results"))
			return
		}

		render.JSON(w, r, testResultsResponse{
			Response: response.OK(),
			Data:     tests,
		})
	}
}
