package res

import (
	"project-go/internal/models"
	"time"
)

type ChatResponse struct {
	ID           uint          `json:"id"`
	SessionID    uint          `json:"session_id"`
	SessionTitle string        `json:"session_title"`
	CreatedAt    time.Time     `json:"created_at"`
	Content      string        `json:"content"`
	Role         string        `json:"role"`
	Status       models.Status `json:"status"`
}

func ChatResponseFromModel(chat *models.ChatMessage) ChatResponse {
	return ChatResponse{
		ID:           chat.ID,
		SessionID:    chat.SessionID,
		SessionTitle: chat.Session.Title,
		CreatedAt:    chat.CreatedAt,
		Content:      chat.Content,
		Role:         chat.Role,
		Status:       chat.Status,
	}
}
