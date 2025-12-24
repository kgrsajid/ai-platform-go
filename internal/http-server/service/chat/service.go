package chatservice

import (
	"errors"
	"project-go/internal/models"
)

type ChatRepository interface {
	CreateChat(chat *models.ChatMessage) (*models.ChatMessage, error)
	GetChatBySessionId(sessionId uint) ([]models.ChatMessage, error)
}

type SessionRepository interface {
	CreateSession(session *models.SessionHistory) (*models.SessionHistory, error)
}
type Service struct {
	chatRepo    ChatRepository
	sessionRepo SessionRepository
}

func New(chatRepo ChatRepository, sessionRepo SessionRepository) *Service {
	return &Service{
		chatRepo:    chatRepo,
		sessionRepo: sessionRepo,
	}
}

func (s *Service) AddMessage(userID uint, sessionID uint, message string) (*models.ChatMessage, error) {
	if message == "" {
		return nil, errors.New("message cannot be empty")
	}
	chat := &models.ChatMessage{
		SessionID: sessionID,
		Role:      "user",
		Content:   message,
	}
	return s.chatRepo.CreateChat(chat)
}

func (s *Service) AddMessageByCreatingSession(userID uint, message string) (*models.ChatMessage, error) {
	if message == "" {
		return nil, errors.New("message cannot be empty")
	}
	session := &models.SessionHistory{
		StudentID: userID,
		Title:     "Dragon history",
	}
	newSession, err := s.sessionRepo.CreateSession(session)
	if err != nil {
		return nil, err
	}
	var sessionID = &newSession.ID
	chat := &models.ChatMessage{
		SessionID: *sessionID,
		Role:      "user",
		Content:   message,
	}
	return s.chatRepo.CreateChat(chat)
}

func (s *Service) GetChatBySessionId(sessionId uint) ([]models.ChatMessage, error) {
	chats, err := s.chatRepo.GetChatBySessionId(sessionId)
	if err != nil {
		return nil, err
	}
	return chats, nil
}
