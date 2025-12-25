package addResult

import (
	"log/slog"
	"net/http"
	req "project-go/internal/http-server/dto/request"
	res "project-go/internal/http-server/dto/response"
	testservice "project-go/internal/http-server/service/test"
	"project-go/internal/lib/api/response"
	"project-go/internal/lib/auth"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Response struct {
	response.Response
	TestResult *res.TestResultResponse
}

func New(log *slog.Logger, testService *testservice.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.test.addResult.New"
		log = log.With(
			slog.String("op", op),
			slog.String("Req_id", middleware.GetReqID(r.Context())),
		)
		var testReq req.TestResultReq
		err := render.DecodeJSON(r.Body, &testReq)
		if err != nil {
			log.Error("failed to decode JSON")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("failed to decode JSON"))
			return
		}
		if testReq.Score < 0 || testReq.MaxScore < 0 || testReq.TestId == 0 {
			log.Error("invalid request")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("invalid request"))
			return
		}
		if testReq.Score > testReq.MaxScore {
			log.Error("Score can't be more than max score")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("Score can't be more than max score"))
			return
		}
		userId, ok := auth.GetUserID(r)
		if !ok {
			log.Error("user id is null")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("user id is null"))
			return
		}
		testReq.UserId = userId
		createdTest, err := testService.AddTestResult(testReq)
		if err != nil {
			log.Error("failed to add result")
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to add result"))
			return
		}

		render.JSON(w, r, Response{
			Response:   response.OK(),
			TestResult: createdTest,
		})

	}
}
