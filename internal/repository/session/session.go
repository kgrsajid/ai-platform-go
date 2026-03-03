package session

import (
	"project-go/internal/models"

	"gorm.io/gorm"
)

type SessionRepository struct {
	db *gorm.DB
}

func NewSessionRepo(db *gorm.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

func (r *SessionRepository) CreateSession(session *models.SessionHistory) (*models.SessionHistory, error) {
	if err := r.db.Create(&session).Error; err != nil {
		return nil, err
	}
	if err := r.db.Preload("Chats").First(&session, session.ID).Error; err != nil {
		return nil, err
	}
	return session, nil
}

func (r *SessionRepository) GetAllSessions(userId uint) ([]models.SessionHistory, error) {
	var sessions []models.SessionHistory
	err := r.db.
		Where("student_id = ?", userId).
		Order("created_at DESC").
		Find(&sessions).Error
	if err != nil {
		return nil, err
	}
	return sessions, nil
}
