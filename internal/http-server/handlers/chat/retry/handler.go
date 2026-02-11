package handler

import (
	"log/slog"
	"net/http"
	res "project-go/internal/http-server/dto/response"
	chatservice "project-go/internal/http-server/service/chat"
	"project-go/internal/lib/api/response"
	"project-go/internal/lib/auth"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Response struct {
	response.Response
	Chat res.ChatResponse `json:"chat"`
}

func New(log *slog.Logger, service *chatservice.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.chat.retry.New"
		log = log.With(
			slog.String("op", op),
			slog.String("req_id", middleware.GetReqID(r.Context())),
		)
		sessionIdStr := chi.URLParam(r, "sessionId")
		if sessionIdStr == "" {
			log.Error("session id is null")
			render.JSON(w, r, response.Error("session id is null"))
			return
		}
		sessionId, err := strconv.ParseUint(sessionIdStr, 10, 64)

		userID, ok := auth.GetUserID(r)

		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		if err != nil {
			log.Error("invalid session id")
			render.JSON(w, r, response.Error("invalid session id"))
			return
		}
		ctx := r.Context()
		message, err := service.RetryLastMessage(ctx, userID, uint(sessionId))
		if err != nil {
			log.Error("internal server error")
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("internal server error"))
			return
		}

		render.JSON(w, r, res.ChatResponseFromModel(message))
	}
}
