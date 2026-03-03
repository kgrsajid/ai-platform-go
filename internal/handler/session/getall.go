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

type sessionListResponse struct {
	response.Response
	Sessions []res.SessionResponse `json:"sessions"`
}

func GetAll(log *slog.Logger, svc *sessionservice.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.session.GetAll"
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

		sessions, err := svc.GetAllSessions(userId)
		if err != nil {
			log.Error("failed to get sessions", slog.String("error", err.Error()))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to get all sessions"))
			return
		}

		responses := make([]res.SessionResponse, len(sessions))
		for i, s := range sessions {
			responses[i] = res.SessionResponse{
				ID:        s.ID,
				Title:     s.Title,
				CreatedAt: s.CreatedAt,
			}
		}

		render.JSON(w, r, sessionListResponse{
			Response: response.OK(),
			Sessions: responses,
		})
	}
}
