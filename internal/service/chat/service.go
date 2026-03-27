package chatservice

import (
	"context"
	"errors"
	"log/slog"
	client "project-go/internal/client/chat"
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
	UpdateTitle(sessionID uint, title string) error
}

type ProgressionRepository interface {
	AddXP(userID uint, xpAmount int, source string) (*models.UserProgress, error)
	AddPoints(userID uint, amount int, source, referenceID, description string) error
	RecordActivity(userID uint) (*models.DailyActivity, *models.UserProgress, error)
	GetOrCreateProgress(userID uint) (*models.UserProgress, *models.User, error)
}

type Service struct {
	chatRepo    ChatRepository
	sessionRepo SessionRepository
	progRepo    ProgressionRepository
	aiAPI       client.AIClient
	log         *slog.Logger
}

func New(chatRepo ChatRepository, sessionRepo SessionRepository, progRepo ProgressionRepository, aiAPI client.AIClient, log *slog.Logger) *Service {
	return &Service{
		chatRepo:    chatRepo,
		sessionRepo: sessionRepo,
		progRepo:    progRepo,
		aiAPI:       aiAPI,
		log:         log,
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

	grade := s.getUserGrade(userID)
	resp, err := s.aiAPI.SendMessage(ctx, userID, lastUserMsg.Content, "ru", grade)
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

// getUserGrade returns the user's grade, defaulting to 5 if unavailable
func (s *Service) getUserGrade(userID uint) int {
	_, user, err := s.progRepo.GetOrCreateProgress(userID)
	if err != nil || user == nil {
		return 5
	}
	if user.Grade <= 0 || user.Grade > 11 {
		return 5
	}
	return user.Grade
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

	// 2️⃣ Запускаем генерацию тайтла параллельно с ответом ИИ (только для первого сообщения)

	// 3️⃣ получаем ответ бота (параллельно с генерацией тайтла)
	botText := "something went wrong"
	grade := s.getUserGrade(userID)

	if s.aiAPI != nil {
		if summary == 0 {
			resp, err := s.aiAPI.SendMessage(ctx, userID, message, "ru", grade)
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

	// Ждём результат генерации тайтла и сохраняем (он уже готов или скоро будет)

	// 3️⃣ сохраняем сообщение бота
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

	// 🎮 Award XP and points for chat message (fire and forget)
	go s.awardChatGamification(userID)

	return botMsg, nil
}

// awardChatGamification awards XP and points for chat messages
func (s *Service) awardChatGamification(userID uint) {
	// Get user's grade for band calculation
	_, user, err := s.progRepo.GetOrCreateProgress(userID)
	if err != nil {
		s.log.Error("gamification: failed to get user progress", slog.String("error", err.Error()))
		return
	}

	grade := user.Grade
	if grade == 0 {
		grade = 5
	}
	band := models.GetGradeBand(grade)

	// Award XP for chat (scaled by grade band)
	xp := 2
	if band == models.BandSprouts {
		xp = 1
	} else if band == models.BandChampions {
		xp = 3
	}

	if _, err := s.progRepo.AddXP(userID, xp, "chat"); err != nil {
		s.log.Error("gamification: failed to add chat XP", slog.String("error", err.Error()))
	}

	// Award small amount of points
	points := 1
	if err := s.progRepo.AddPoints(userID, points, "chat", "", "Chat message"); err != nil {
		s.log.Error("gamification: failed to add chat points", slog.String("error", err.Error()))
	}

	// Record activity for streak tracking
	if _, _, err := s.progRepo.RecordActivity(userID); err != nil {
		s.log.Error("gamification: failed to record chat activity", slog.String("error", err.Error()))
	}

	s.log.Info("gamification: chat reward awarded",
		slog.Int("userID", int(userID)),
		slog.Int("xp", xp),
		slog.Int("points", points),
	)
}

func (s *Service) AddMessageByCreatingSession(ctx context.Context, userID uint, message string) (*models.ChatMessage, error) {
	if message == "" {
		return nil, errors.New("message cannot be empty")
	}

	title, err := s.aiAPI.GenerateTitle(context.Background(), message, "ru")
	if err != nil || title == "" {
		title = "Новый чат"
	}

	session := &models.SessionHistory{
		StudentID: userID,
		Title:     title,
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
	grade := s.getUserGrade(userID)
	if s.aiAPI != nil {
		resp, err := s.aiAPI.SendMessage(ctx, userID, message, "ru", grade)
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

	// 🎮 Award XP and points for chat message
	go s.awardChatGamification(userID)

	return botMsg, nil
}

func (s *Service) GetChatBySessionId(sessionId uint) ([]models.ChatMessage, error) {
	return s.chatRepo.GetChatBySessionId(sessionId)
}
