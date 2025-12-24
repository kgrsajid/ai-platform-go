package getbyid

import (
	"log/slog"
	"net/http"
	res "project-go/internal/http-server/dto/response"
	testservice "project-go/internal/http-server/service/test"
	"project-go/internal/lib/api/response"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Response struct {
	response.Response
	Test *res.TestDetailsResponse `json:"test"`
}

func New(log *slog.Logger, testService *testservice.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.test.getById.New"
		log = log.With(
			slog.String("op", op),
			slog.String("Req_id", middleware.GetReqID(r.Context())),
		)
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
		test, err := testService.GetTestById(testId)
		if err != nil {
			log.Error("failed to get test by id")
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to get test by id"))
			return
		}
		render.JSON(w, r, Response{
			Response: response.OK(),
			Test:     test,
		})
	}
}
