package chatservice

import (
	"context"
	"errors"
	"project-go/internal/client/chat"
	"project-go/internal/models"
	"strings"

	"gorm.io/gorm"
)

type ChatRepository interface {
	CreateChat(chat *models.ChatMessage) (*models.ChatMessage, error)
	GetChatBySessionId(sessionId uint) ([]models.ChatMessage, error)
	CreateChatTx(tx *gorm.DB, chat *models.ChatMessage) (*models.ChatMessage, error)
	UpdateChat(chat *models.ChatMessage) (*models.ChatMessage, error)
	GetLastErrorUserMessage(sessionID uint) (*models.ChatMessage, error)
	BeginTx() *gorm.DB
}

type SessionRepository interface {
	CreateSession(session *models.SessionHistory) (*models.SessionHistory, error)
}

type Service struct {
	chatRepo    ChatRepository
	sessionRepo SessionRepository
	aiAPI       client.AIClient
}

func New(chatRepo ChatRepository, sessionRepo SessionRepository, aiAPI client.AIClient) *Service {
	return &Service{
		chatRepo:    chatRepo,
		sessionRepo: sessionRepo,
		aiAPI:       aiAPI,
	}
}

func (s *Service) RetryLastMessage(ctx context.Context, userID, sessionID uint) (*models.ChatMessage, error) {
	if s.aiAPI == nil {
		return nil, errors.New("ai api is not configured")
	}

	lastUserMsg, err := s.chatRepo.GetLastErrorUserMessage(sessionID)
	if err != nil {
		return nil, err
	}
	if lastUserMsg == nil {
		return nil, errors.New("no message to retry")
	}

	resp, err := s.aiAPI.SendMessage(ctx, userID, lastUserMsg.Content, "ru")
	if err != nil {
		return nil, err
	}

	botMsg := &models.ChatMessage{
		SessionID: sessionID,
		Role:      "bot",
		Content:   resp.Response,
	}
	if _, err := s.chatRepo.CreateChat(botMsg); err != nil {
		return nil, err
	}

	lastUserMsg.Status = models.Success
	if _, err := s.chatRepo.UpdateChat(lastUserMsg); err != nil {
		return nil, err
	}

	return botMsg, nil
}

func (s *Service) AddMessage(ctx context.Context, userID, sessionID uint, message string, summary int) (*models.ChatMessage, error) {
	if strings.TrimSpace(message) == "" {
		return nil, errors.New("empty message")
	}

	userMsg := &models.ChatMessage{
		SessionID: sessionID,
		Role:      "user",
		Content:   message,
		Status:    models.Pending,
	}

	newUserChat, err := s.chatRepo.CreateChat(userMsg)
	if err != nil {
		return nil, err
	}

	botText := "something went wrong"

	if s.aiAPI != nil {
		if summary == 0 {
			resp, err := s.aiAPI.SendMessage(ctx, userID, message, "ru")
			if err != nil {
				newUserChat.Status = models.Error
				if _, updateErr := s.chatRepo.UpdateChat(newUserChat); updateErr != nil {
					return nil, updateErr
				}
				return nil, err
			}
			botText = resp.Response
		} else {
			resp, err := s.aiAPI.CreateSummary(ctx, userID, message, "ru")
			if err != nil {
				return nil, err
			}
			botText = resp.Summary
		}
	}

	botMsg := &models.ChatMessage{
		SessionID: sessionID,
		Role:      "bot",
		Content:   botText,
	}
	if _, err := s.chatRepo.CreateChat(botMsg); err != nil {
		return nil, err
	}

	newUserChat.Status = models.Success
	if _, updateErr := s.chatRepo.UpdateChat(newUserChat); updateErr != nil {
		return nil, updateErr
	}

	return botMsg, nil
}

func (s *Service) AddMessageByCreatingSession(ctx context.Context, userID uint, message string) (*models.ChatMessage, error) {
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

	chat := &models.ChatMessage{
		SessionID: newSession.ID,
		Role:      "user",
		Content:   message,
		Status:    models.Pending,
	}
	newUserChat, err := s.chatRepo.CreateChat(chat)
	if err != nil {
		return nil, err
	}

	botText := "something went wrong"
	if s.aiAPI != nil {
		resp, err := s.aiAPI.SendMessage(ctx, userID, message, "ru")
		if err != nil {
			newUserChat.Status = models.Error
			if _, updateErr := s.chatRepo.UpdateChat(newUserChat); updateErr != nil {
				return nil, updateErr
			}
			return nil, err
		}
		botText = resp.Response
	}

	botMsg := &models.ChatMessage{
		SessionID: newSession.ID,
		Role:      "bot",
		Content:   botText,
	}
	if _, err := s.chatRepo.CreateChat(botMsg); err != nil {
		return nil, err
	}

	newUserChat.Status = models.Success
	if _, updateErr := s.chatRepo.UpdateChat(newUserChat); updateErr != nil {
		return nil, updateErr
	}

	return botMsg, nil
}

func (s *Service) GetChatBySessionId(sessionId uint) ([]models.ChatMessage, error) {
	return s.chatRepo.GetChatBySessionId(sessionId)
}
