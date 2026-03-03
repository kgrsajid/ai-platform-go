package chat

import (
	"errors"
	"project-go/internal/models"

	"gorm.io/gorm"
)

type ChatRepository struct {
	db *gorm.DB
}

func NewChatRepo(db *gorm.DB) *ChatRepository {
	return &ChatRepository{db: db}
}

func (r *ChatRepository) BeginTx() *gorm.DB {
	return r.db.Begin()
}

func (r *ChatRepository) CreateChatTx(tx *gorm.DB, chat *models.ChatMessage) (*models.ChatMessage, error) {
	if err := tx.Create(chat).Error; err != nil {
		return nil, err
	}

	if err := tx.
		Preload("Session").
		Preload("Session.Student").
		First(chat, chat.ID).Error; err != nil {
		return nil, err
	}

	return chat, nil
}

func (r *ChatRepository) UpdateChat(chat *models.ChatMessage) (*models.ChatMessage, error) {
	if err := r.db.
		Model(&models.ChatMessage{}).
		Where("id = ?", chat.ID).
		Updates(chat).Error; err != nil {
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

func (r *ChatRepository) GetLastErrorUserMessage(sessionID uint) (*models.ChatMessage, error) {
	var msg models.ChatMessage

	err := r.db.
		Where("session_id = ? AND role = ? AND status = ?", sessionID, "user", models.Error).
		Order("created_at DESC").
		First(&msg).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &msg, nil
}
