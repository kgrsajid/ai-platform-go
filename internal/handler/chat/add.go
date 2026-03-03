package chathandler

import (
	"log/slog"
	"net/http"

	"project-go/internal/dto/request"
	res "project-go/internal/dto/response"
	"project-go/internal/lib/auth"
	"project-go/internal/lib/response"
	chatservice "project-go/internal/service/chat"

	"github.com/go-chi/render"
)

type chatResponse struct {
	response.Response
	Chat res.ChatResponse `json:"chat"`
}

func Add(log *slog.Logger, svc *chatservice.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req request.AddChatRequest
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

		chat, err := svc.AddMessage(r.Context(), userID, req.SessionId, req.Message, 0)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error(err.Error()))
			return
		}

		render.JSON(w, r, res.ChatResponseFromModel(chat))
	}
}
