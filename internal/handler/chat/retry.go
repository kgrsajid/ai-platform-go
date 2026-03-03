package chathandler

import (
	"log/slog"
	"net/http"
	"strconv"

	res "project-go/internal/dto/response"
	"project-go/internal/lib/auth"
	"project-go/internal/lib/response"
	chatservice "project-go/internal/service/chat"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

func Retry(log *slog.Logger, svc *chatservice.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.chat.Retry"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		sessionIdStr := chi.URLParam(r, "sessionId")
		if sessionIdStr == "" {
			log.Error("missing session id")
			render.JSON(w, r, response.Error("missing session id"))
			return
		}

		sessionId, err := strconv.ParseUint(sessionIdStr, 10, 64)
		if err != nil {
			log.Error("invalid session id")
			render.JSON(w, r, response.Error("invalid session id"))
			return
		}

		userID, ok := auth.GetUserID(r)
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		message, err := svc.RetryLastMessage(r.Context(), userID, uint(sessionId))
		if err != nil {
			log.Error("failed to retry message", slog.String("error", err.Error()))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("internal server error"))
			return
		}

		render.JSON(w, r, res.ChatResponseFromModel(message))
	}
}
