package getbysessionid

import (
	"log/slog"
	"net/http"
	res "project-go/internal/http-server/dto/response"
	chatservice "project-go/internal/http-server/service/chat"
	"project-go/internal/lib/api/response"
	"project-go/internal/models"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type ChatGetBySessionId interface {
	GetChatBySessionId(sessionId uint) ([]models.ChatMessage, error)
}

type Response struct {
	response.Response
	Chats []res.ChatResponse `json:"chat"`
}

func New(log *slog.Logger, service *chatservice.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.chat.getBySessionId.New"
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
		if err != nil {
			log.Error("invalid session id")
			render.JSON(w, r, response.Error("invalid session id"))
			return
		}
		chats, err := service.GetChatBySessionId(uint(sessionId))

		if err != nil {
			log.Error("failed to get chats by session id")
			render.JSON(w, r, response.Error("failed to get chats by session id"))
			return
		}
		responses := make([]res.ChatResponse, len(chats))
		for i, ch := range chats {
			responses[i] = res.ChatResponseFromModel(&ch)
		}

		render.JSON(w, r, Response{
			Response: response.OK(),
			Chats:    responses,
		})
	}
}
