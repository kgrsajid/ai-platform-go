package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"project-go/internal/dto/request"
	res "project-go/internal/dto/response"
)

type AIClient interface {
	SendMessage(ctx context.Context, userID uint, message string, language string) (*res.AIResponse, error)
	CreateSummary(ctx context.Context, userID uint, topic string, language string) (*res.SummaryResponse, error)
	GenerateTitle(ctx context.Context, message string, language string) (string, error)
	GenerateQuiz(ctx context.Context, quiz request.GenerateQuizReq, language string) (*res.GeneratedTestResponse, error)
	GenerateCards(ctx context.Context, cardPayload request.GenerateCardReq, language string) (*res.GeneratedCardHolderDetailResponse, error)
}

type aiClient struct {
	baseURL    string
	client     *http.Client
	slowClient *http.Client
}

func NewAIClient(baseURL string) AIClient {
	return &aiClient{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
		slowClient: &http.Client{
			Timeout: 5 * time.Minute,
		},
	}
}

const (
	messageEndpoint = "/chat"
	summaryEndpoint = "/summary"
	titleEndpoint   = "/generate-title"
	generateQuiz    = "/quiz/generate-for-platform"
	generateCard    = "/flashcards/generate-for-platform"
)

func (a *aiClient) SendMessage(ctx context.Context, userID uint, message string, language string) (*res.AIResponse, error) {
	body, err := json.Marshal(request.AiRequest{
		UserID:   strconv.FormatUint(uint64(userID), 10),
		Message:  message,
		Language: language,
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, a.baseURL+messageEndpoint, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("ai service returned status " + resp.Status)
	}

	var result res.AIResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

type quizPayload struct {
	Title             string `json:"title"`
	Description       string `json:"context"`
	Categories        []uint `json:"categories"`
	Difficulty        string `json:"difficulty"`
	IsPrivate         bool   `json:"is_private"`
	NumberOfQuestions int    `json:"num_questions"`
	Language          string `json:"language"`
}

func (a *aiClient) GenerateQuiz(ctx context.Context, quiz request.GenerateQuizReq, language string) (*res.GeneratedTestResponse, error) {
	lang := language
	if quiz.Language != nil && *quiz.Language != "" {
		lang = *quiz.Language
	}
	payload := quizPayload{
		Title:             quiz.Title,
		Description:       quiz.Description,
		Categories:        quiz.Categories,
		Difficulty:        quiz.Difficulty,
		IsPrivate:         quiz.IsPrivate,
		NumberOfQuestions: quiz.NumberOfQuestions,
		Language:          lang,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, a.baseURL+generateQuiz, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := a.slowClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("ai service returned status " + resp.Status)
	}

	var raw res.GeneratedTestResponse
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}
	return &raw, nil
}

func (a *aiClient) GenerateCards(ctx context.Context, cardPayload request.GenerateCardReq, language string) (*res.GeneratedCardHolderDetailResponse, error) {
	lang := language
	if cardPayload.Language != nil && *cardPayload.Language != "" {
		lang = *cardPayload.Language
	}
	cardPayload.Language = &lang
	body, err := json.Marshal(cardPayload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, a.baseURL+generateCard, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := a.slowClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("ai service returned status " + resp.Status)
	}
	var raw res.GeneratedCardHolderDetailResponse
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}
	return &raw, nil
}

func (a *aiClient) GenerateTitle(ctx context.Context, message string, language string) (string, error) {
	body, err := json.Marshal(request.TitleRequest{
		Message:  message,
		Language: language,
	})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, a.baseURL+titleEndpoint, bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("ai service returned status " + resp.Status)
	}

	var result res.TitleResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.Title, nil
}

func (a *aiClient) CreateSummary(ctx context.Context, userID uint, topic string, language string) (*res.SummaryResponse, error) {
	body, err := json.Marshal(request.SummaryRequest{
		UserID:   strconv.FormatUint(uint64(userID), 10),
		Topic:    topic,
		Language: language,
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, a.baseURL+summaryEndpoint, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("ai service returned status " + resp.Status)
	}

	var result res.SummaryResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}
