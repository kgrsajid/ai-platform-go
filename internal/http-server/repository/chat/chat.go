package chat

import (
	"project-go/internal/models"

	"gorm.io/gorm"
)

type ChatRepository struct {
	db *gorm.DB
}

func NewChatnRepo(db *gorm.DB) *ChatRepository {
	return &ChatRepository{db: db}
}

func (r *ChatRepository) CreateChat(chat *models.ChatMessage) (*models.ChatMessage, error) {
	if err := r.db.Create(chat).Error; err != nil {
		return nil, err
	}
	if err := r.db.
		Preload("Session").
		Preload("Session.Student").
		First(chat, chat.ID).Error; err != nil {
		return nil, err
	}

	return chat, nil
}

func (r *ChatRepository) GetChatBySessionId(sessionId uint) ([]models.ChatMessage, error) {
	var chats []models.ChatMessage
	if err := r.db.
		Preload("Session").
		Preload("Session.Student").
		Where("session_id = ?", sessionId).
		Order("created_at ASC").
		Find(&chats).Error; err != nil {
		return nil, err
	}
	return chats, nil
}
