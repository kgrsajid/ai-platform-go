package getAllSession

import (
	"log/slog"
	"net/http"
	res "project-go/internal/http-server/dto/response"
	sessionService "project-go/internal/http-server/service/session"
	"project-go/internal/lib/api/response"
	"project-go/internal/lib/auth"
	"project-go/internal/models"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type SessionGetAll interface {
	GetAllSessions(userId uint) ([]models.SessionHistory, error)
}
type Response struct {
	response.Response
	Sessions []res.SessionResponse `json:"sessions"`
}

func New(log *slog.Logger, sessionService *sessionService.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.session.getAll.New"
		log = log.With(
			slog.String("op", op),
			slog.String("req_id", middleware.GetReqID(r.Context())),
		)
		userId, ok := auth.GetUserID(r)
		if !ok {
			log.Error("user id is null")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("user id is null"))
			return
		}
		sessions, err := sessionService.GetAllSessions(userId)
		if err != nil {
			log.Error("failed to get all sessions")
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
		render.JSON(w, r, Response{
			Response: response.OK(),
			Sessions: responses,
		})
	}
}
