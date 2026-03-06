package sessionservice

import (
	"context"
	client "project-go/internal/client/chat"
	"project-go/internal/models"
)

type SessionRepository interface {
	GetAllSessions(userId uint) ([]models.SessionHistory, error)
	CreateSession(session *models.SessionHistory) (*models.SessionHistory, error)
	DeleteSession(sessionID uint) error
}

type Service struct {
	sessionRepo SessionRepository
	aiAPI       client.AIClient
}

func New(sessionRepo SessionRepository, aiClient client.AIClient) *Service {
	return &Service{sessionRepo: sessionRepo, aiAPI: aiClient}
}

func (s *Service) GetAllSessions(userId uint) ([]models.SessionHistory, error) {
	return s.sessionRepo.GetAllSessions(userId)
}

func (s *Service) DeleteSession(sessionID uint) error {
	return s.sessionRepo.DeleteSession(sessionID)
}

func (s *Service) CreateSession(userID uint, message string) (*models.SessionHistory, error) {
	title, err := s.aiAPI.GenerateTitle(context.Background(), message, "ru")
	if err != nil || title == "" {
		title = "Новый чат"
	}
	session := &models.SessionHistory{
		StudentID: userID,
		Title:     title,
	}
	return s.sessionRepo.CreateSession(session)
}
