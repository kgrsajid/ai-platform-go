package sessionhandler

import (
	"log/slog"
	"net/http"

	res "project-go/internal/dto/response"
	"project-go/internal/lib/auth"
	"project-go/internal/lib/response"
	sessionservice "project-go/internal/service/session"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type sessionResponse struct {
	response.Response
	Data res.SessionResponse `json:"data"`
}

func Create(log *slog.Logger, svc *sessionservice.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.session.Create"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		userId, ok := auth.GetUserID(r)
		if !ok {
			log.Error("unauthorized: missing user id")
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, response.Error("unauthorized"))
			return
		}

		newSession, err := svc.CreateSession(userId)
		if err != nil {
			log.Error("failed to create session", slog.String("error", err.Error()))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("internal server error"))
			return
		}

		render.JSON(w, r, sessionResponse{
			Response: response.OK(),
			Data: res.SessionResponse{
				ID:        newSession.ID,
				Title:     newSession.Title,
				CreatedAt: newSession.CreatedAt,
			},
		})
	}
}
