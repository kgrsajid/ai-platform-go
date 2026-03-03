package chathandler

import (
	"log/slog"
	"net/http"
	"strconv"

	res "project-go/internal/dto/response"
	"project-go/internal/lib/response"
	chatservice "project-go/internal/service/chat"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type chatListResponse struct {
	response.Response
	Chats []res.ChatResponse `json:"chat"`
}

func GetBySessionID(log *slog.Logger, svc *chatservice.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.chat.GetBySessionID"
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

		chats, err := svc.GetChatBySessionId(uint(sessionId))
		if err != nil {
			log.Error("failed to get chats", slog.String("error", err.Error()))
			render.JSON(w, r, response.Error("failed to get chats by session id"))
			return
		}

		responses := make([]res.ChatResponse, len(chats))
		for i, ch := range chats {
			responses[i] = res.ChatResponseFromModel(&ch)
		}

		render.JSON(w, r, chatListResponse{
			Response: response.OK(),
			Chats:    responses,
		})
	}
}
