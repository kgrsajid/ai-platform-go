package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	req "project-go/internal/http-server/dto/request"
	res "project-go/internal/http-server/dto/response"
	"strconv"
	"time"
)

type AIClient interface {
	SendMessage(
		ctx context.Context,
		userID uint,
		message string,
		language string,
	) (*res.AIResponse, error)

	CreateSummary(
		ctx context.Context,
		userID uint,
		topic string,
		language string,
	) (*res.SummaryResponse, error)
}

type aiClient struct {
	baseURL string
	client  *http.Client
}

func NewAIClient(baseURL string) AIClient {
	return &aiClient{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

const (
	messageEndpoint = "/chat"
	summaryEndpoint = "/summary"
)

func (a *aiClient) SendMessage(
	ctx context.Context,
	userID uint,
	message string,
	language string,
) (*res.AIResponse, error) {

	body, err := json.Marshal(req.AiRequest{
		UserID:   string(userID),
		Message:  message,
		Language: language,
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		a.baseURL+messageEndpoint,
		bytes.NewBuffer(body),
	)
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

	var res res.AIResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (a *aiClient) CreateSummary(
	ctx context.Context,
	userID uint,
	topic string,
	language string,
) (*res.SummaryResponse, error) {

	body, err := json.Marshal(req.SummaryRequest{
		UserID:   strconv.FormatUint(uint64(userID), 10),
		Topic:    topic,
		Language: language,
	})
	if err != nil {
		return nil, err
	}

	reqHTTP, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		a.baseURL+summaryEndpoint,
		bytes.NewBuffer(body),
	)
	if err != nil {
		return nil, err
	}

	reqHTTP.Header.Set("Content-Type", "application/json")

	resp, err := a.client.Do(reqHTTP)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("ai service returned status " + resp.Status)
	}

	var response res.SummaryResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}
