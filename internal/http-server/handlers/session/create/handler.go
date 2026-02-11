package create

import (
	"log/slog"
	"net/http"
	res "project-go/internal/http-server/dto/response"
	sessionService "project-go/internal/http-server/service/session"
	"project-go/internal/lib/api/response"
	"project-go/internal/lib/auth"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Response struct {
	response.Response
	Data res.SessionResponse `json:"data"`
}

func New(log *slog.Logger, chatService *sessionService.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.session.create.New"
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

		newSession, err := chatService.CreateSession(userId)
		if err != nil {
			log.Error("internal server error")
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("internal server error"))
			return
		}

		render.JSON(w, r, Response{
			Response: response.OK(),
			Data: res.SessionResponse{
				ID:        newSession.ID,
				Title:     newSession.Title,
				CreatedAt: newSession.CreatedAt,
			},
		})
	}
}
