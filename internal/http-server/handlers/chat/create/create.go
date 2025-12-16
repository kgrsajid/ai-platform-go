package chatCreate

import (
	"log/slog"
	"net/http"
	req "project-go/internal/http-server/dto/request"
	res "project-go/internal/http-server/dto/response"
	chatservice "project-go/internal/http-server/service/chat"
	"project-go/internal/lib/api/response"
	"project-go/internal/lib/auth"
	"project-go/internal/models"

	"github.com/go-chi/render"
)

type ChatCreate interface {
	CreateChat(chat *models.ChatMessage) (*models.ChatMessage, error)
}

type SessionCreate interface {
	CreateSession(session *models.SessionHistory) (*models.SessionHistory, error)
}

type Response struct {
	response.Response
	Chat res.ChatResponse
}

func New(log *slog.Logger, service *chatservice.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var req req.CreateChatRequest
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("invalid request body"))
			return
		}

		userID, ok := auth.GetUserID(r)
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		chat, err := service.CreateMessage(userID, req.SessionId, req.Message)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error(err.Error()))
			return
		}

		render.JSON(w, r, res.ChatResponseFromModel(chat))
	}
}
