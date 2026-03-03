package sessionhandler

import (
	"log/slog"
	"net/http"
	"strconv"

	"project-go/internal/lib/response"
	sessionservice "project-go/internal/service/session"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

func Delete(log *slog.Logger, svc *sessionservice.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.session.Delete"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		sessionIdStr := chi.URLParam(r, "sessionId")
		sessionId, err := strconv.ParseUint(sessionIdStr, 10, 64)
		if err != nil {
			log.Error("invalid session id", slog.String("error", err.Error()))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("invalid session id"))
			return
		}

		if err := svc.DeleteSession(uint(sessionId)); err != nil {
			log.Error("failed to delete session", slog.String("error", err.Error()))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to delete session"))
			return
		}

		render.JSON(w, r, response.OK())
	}
}
