package testhandler

import (
	"log/slog"
	"net/http"

	"project-go/internal/dto/request"
	"project-go/internal/lib/auth"
	"project-go/internal/lib/response"
	testviewservice "project-go/internal/service/testview"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

func AddView(log *slog.Logger, svc *testviewservice.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.test.AddView"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req request.TestViewReq
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("failed to decode request")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("invalid request body"))
			return
		}

		userId, ok := auth.GetUserID(r)
		if !ok {
			log.Error("unauthorized: missing user id")
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, response.Error("unauthorized"))
			return
		}

		if req.TestId == 0 {
			log.Error("missing test id")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("test id is required"))
			return
		}

		if err := svc.AddTestView(req.TestId, userId); err != nil {
			log.Error("failed to add test view", slog.String("error", err.Error()))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("internal server error"))
			return
		}

		render.JSON(w, r, response.OK())
	}
}
