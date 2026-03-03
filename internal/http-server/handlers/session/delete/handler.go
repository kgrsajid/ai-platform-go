package delete

import (
	"log/slog"
	"net/http"
	"strconv"

	"project-go/internal/lib/api/response"
	sessionService "project-go/internal/http-server/service/session"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

func New(log *slog.Logger, svc *sessionService.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.session.delete.New"
		log = log.With(
			slog.String("op", op),
			slog.String("req_id", middleware.GetReqID(r.Context())),
		)

		sessionIDStr := chi.URLParam(r, "sessionId")
		sessionID, err := strconv.ParseUint(sessionIDStr, 10, 64)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("invalid session id"))
			return
		}

		if err := svc.DeleteSession(uint(sessionID)); err != nil {
			log.Error("failed to delete session", slog.String("error", err.Error()))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to delete session"))
			return
		}

		render.JSON(w, r, response.OK())
	}
}
