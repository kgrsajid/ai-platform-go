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

type testResultResponse struct {
	response.Response
	TestResult *res.TestResultResponse `json:"data"`
}

func AddResult(log *slog.Logger, svc *testservice.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.test.AddResult"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var testReq request.TestResultReq
		if err := render.DecodeJSON(r.Body, &testReq); err != nil {
			log.Error("failed to decode request")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("invalid request body"))
			return
		}

		if testReq.Score < 0 || testReq.MaxScore < 0 || testReq.TestId == 0 {
			log.Error("invalid request fields")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("invalid request"))
			return
		}

		if testReq.Score > testReq.MaxScore {
			log.Error("score exceeds max score")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("score cannot exceed max score"))
			return
		}

		userId, ok := auth.GetUserID(r)
		if !ok {
			log.Error("unauthorized: missing user id")
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, response.Error("unauthorized"))
			return
		}
		testReq.UserId = userId

		result, err := svc.AddTestResult(testReq)
		if err != nil {
			log.Error("failed to add result", slog.String("error", err.Error()))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to add result"))
			return
		}

		render.JSON(w, r, testResultResponse{
			Response:   response.OK(),
			TestResult: result,
		})
	}
}
