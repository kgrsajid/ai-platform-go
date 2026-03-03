package chatservice

import (
	"context"
	"errors"
	"fmt"
	client "project-go/internal/http-server/client/chat"
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

type Service struct {
	chatRepo    ChatRepository
	sessionRepo SessionRepository
	aiAPi       client.AIClient
}

func New(chatRepo ChatRepository, sessionRepo SessionRepository, aiApi client.AIClient) *Service {
	return &Service{
		chatRepo:    chatRepo,
		sessionRepo: sessionRepo,
		aiAPi:       aiApi,
	}
}

func (s *Service) RetryLastMessage(
	ctx context.Context,
	userID, sessionID uint,
) (*models.ChatMessage, error) {

	if s.aiAPi == nil {
		return nil, errors.New("ai api is not configured")
	}

	// 1️⃣ получаем последнее упавшее сообщение пользователя
	lastUserMsg, err := s.chatRepo.GetLastErrorUserMessage(sessionID)
	if err != nil {
		return nil, err
	}

	if lastUserMsg == nil {
		return nil, errors.New("no message to retry")
	}

	// 2️⃣ повторно отправляем в AI
	resp, err := s.aiAPi.SendMessage(
		ctx,
		userID,
		lastUserMsg.Content,
		"ru",
	)

	if err != nil {
		// всё ещё ошибка — статус не меняем
		return nil, err
	}

	// 3️⃣ сохраняем сообщение бота
	botMsg := &models.ChatMessage{
		SessionID: sessionID,
		Role:      "bot",
		Content:   resp.Response,
	}

	if _, err := s.chatRepo.CreateChat(botMsg); err != nil {
		return nil, err
	}

	// 4️⃣ обновляем статус user-сообщения
	lastUserMsg.Status = models.Success
	if _, err := s.chatRepo.UpdateChat(lastUserMsg); err != nil {
		return nil, err
	}

	return botMsg, nil
}

func (s *Service) AddMessage(
	ctx context.Context,
	userID, sessionID uint,
	message string,
	summary int,
) (*models.ChatMessage, error) {

	if strings.TrimSpace(message) == "" {
		return nil, errors.New("empty message")
	}

	// 1️⃣ сохраняем сообщение пользователя
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
	titleCh := make(chan string, 1)
	isFirstMessage := false
	if existing, countErr := s.chatRepo.GetChatBySessionId(sessionID); countErr == nil && len(existing) == 1 && s.aiAPi != nil {
		isFirstMessage = true
		go func() {
			title, err := s.aiAPi.GenerateTitle(context.Background(), message, "ru")
			if err == nil && title != "" {
				titleCh <- title
			} else {
				titleCh <- ""
			}
		}()
	}

	// 3️⃣ получаем ответ бота (параллельно с генерацией тайтла)
	botText := "something went wrong"

	if s.aiAPi != nil {
		if summary == 0 {
			fmt.Println("chat is working now")
			resp, err := s.aiAPi.SendMessage(ctx, userID, message, "ru")
			if err != nil {
				newUserChat.Status = models.Error
				if _, updateErr := s.chatRepo.UpdateChat(newUserChat); updateErr != nil {
					return nil, updateErr
				}
				return nil, err
			}
			botText = resp.Response
		} else {

			fmt.Println("summary is working now")
			resp, err := s.aiAPi.CreateSummary(ctx, userID, message, "ru")
			if err != nil {
				return nil, err
			}
			botText = resp.Summary
		}
	}

	// Ждём результат генерации тайтла и сохраняем (он уже готов или скоро будет)
	if isFirstMessage {
		if title := <-titleCh; title != "" {
			_ = s.sessionRepo.UpdateTitle(sessionID, title)
		}
	}

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

	// 🔥 ВОЗВРАЩАЕМ БОТА
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
	var sessionID = &newSession.ID
	chat := &models.ChatMessage{
		SessionID: *sessionID,
		Role:      "user",
		Content:   message,
		Status:    models.Pending,
	}
	newUserChat, err := s.chatRepo.CreateChat(chat)
	if err != nil {
		return nil, err
	}
	botText := "something went wrong"

	if s.aiAPi != nil {
		resp, err := s.aiAPi.SendMessage(ctx, userID, message, "ru")
		if err != nil {
			newUserChat.Status = models.Error
			if _, updateErr := s.chatRepo.UpdateChat(newUserChat); updateErr != nil {
				return nil, updateErr
			}
			return nil, err
		}

		botText = resp.Response
	}

	// 3️⃣ сохраняем сообщение бота
	botMsg := &models.ChatMessage{
		SessionID: *sessionID,
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

	// 🔥 ВОЗВРАЩАЕМ БОТА
	return botMsg, nil
}

func (s *Service) GetChatBySessionId(sessionId uint) ([]models.ChatMessage, error) {
	chats, err := s.chatRepo.GetChatBySessionId(sessionId)
	if err != nil {
		return nil, err
	}
	return chats, nil
}
