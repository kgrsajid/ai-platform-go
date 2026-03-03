package sessionservice

import "project-go/internal/models"

type SessionRepository interface {
	GetAllSessions(userId uint) ([]models.SessionHistory, error)
	CreateSession(session *models.SessionHistory) (*models.SessionHistory, error)
}

type Service struct {
	sessionRepo SessionRepository
}

func New(sessionRepo SessionRepository) *Service {
	return &Service{sessionRepo: sessionRepo}
}

func (s *Service) GetAllSessions(userId uint) ([]models.SessionHistory, error) {
	return s.sessionRepo.GetAllSessions(userId)
}

func (s *Service) CreateSession(userID uint) (*models.SessionHistory, error) {
	session := &models.SessionHistory{
		StudentID: userID,
		Title:     "Dragon history",
	}
	return s.sessionRepo.CreateSession(session)
}
