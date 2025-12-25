package getalluserresults

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
	"github.com/go-chi/render"
)

type Response struct {
	response.Response
	Data []res.TestResultResponse `json:"data"`
}

func New(log *slog.Logger, testService *testservice.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId, ok := auth.GetUserID(r)
		if !ok {
			log.Error("user id is null")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("user id is null"))
			return
		}
		testIdStr := chi.URLParam(r, "testId") // получаем строку из URL
		testId, err := strconv.ParseUint(testIdStr, 10, 64)
		if err != nil {
			http.Error(w, "invalid testId", http.StatusBadRequest)
			return
		}
		duration := r.URL.Query().Get("duration")
		if duration == "" {
			duration = "month"
		}
		uid := uint(userId)
		tid := uint(testId)

		filter := req.GetALlTestResultsFilter{
			UserID:   uid,
			TestID:   tid,
			Duration: duration,
		}

		tests, err := testService.GetAllUserTestResults(filter)
		if err != nil {
			log.Error("failed to get test results")
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to get test results"))
		}

		render.JSON(w, r, Response{
			Response: response.OK(),
			Data:     tests,
		})

	}
}
