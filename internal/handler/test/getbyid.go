package testhandler

import (
	"log/slog"
	"net/http"
	"strconv"

	res "project-go/internal/dto/response"
	"project-go/internal/lib/response"
	testservice "project-go/internal/service/test"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type testByIDResponse struct {
	response.Response
	Test *res.TestDetailsResponse `json:"test"`
}

func GetByID(log *slog.Logger, svc *testservice.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.test.GetByID"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

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

		test, err := svc.GetTestById(testId)
		if err != nil {
			log.Error("failed to get test", slog.String("error", err.Error()))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to get test by id"))
			return
		}

		render.JSON(w, r, testByIDResponse{
			Response: response.OK(),
			Test:     test,
		})
	}
}
