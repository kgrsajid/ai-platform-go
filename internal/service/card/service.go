package cardservice

import (
	"project-go/internal/dto/request"
	res "project-go/internal/dto/response"
	"project-go/internal/models"
)

type CardRepository interface {
	CreateCardHolder(card *models.CardHolder) (*models.CardHolder, error)
	GetAll(filter *request.CardFilter) ([]models.CardHolder, int64, error)
	GetCardsByCardHolderId(cardHolderId uint) (*models.CardHolder, error)
	UpdateCard(cardholder *models.CardHolder) (*models.CardHolder, error)
}

type Service struct {
	cardRepo CardRepository
}

func New(cardRepo CardRepository) *Service {
	return &Service{cardRepo: cardRepo}
}

func (s *Service) CreateCardHolder(req request.CardHolderRequest) (*res.CardHolderDetailResponse, error) {
	modelsCard := toCardHolderModel(req)
	createdCardHolder, err := s.cardRepo.CreateCardHolder(&modelsCard)
	if err != nil {
		return nil, err
	}
	return toCardHolderDetailResponse(createdCardHolder), nil
}

func (s *Service) UpdateCardHolder(req request.CardHolderRequest, cardId uint) (*res.CardHolderDetailResponse, error) {
	modelsCard := toCardHolderModel(req)
	modelsCard.ID = cardId
	updatedCardHolder, err := s.cardRepo.UpdateCard(&modelsCard)
	if err != nil {
		return nil, err
	}
	return toCardHolderDetailResponse(updatedCardHolder), nil
}

func (s *Service) GetAll(cardFilter request.CardFilter) ([]res.CardHolderResponse, int64, error) {
	cardHolderModels, total, err := s.cardRepo.GetAll(&cardFilter)
	if err != nil {
		return nil, 0, err
	}
	return toCardHolderResponse(cardHolderModels), total, nil
}

func (s *Service) GetCardsByCardHolderId(id uint) (*res.CardHolderDetailResponse, error) {
	cardHolder, err := s.cardRepo.GetCardsByCardHolderId(id)
	if err != nil {
		return nil, err
	}
	return toCardHolderDetailResponse(cardHolder), nil
}

func toCardHolderResponse(cardHolderModels []models.CardHolder) []res.CardHolderResponse {
	cardHolders := make([]res.CardHolderResponse, 0, len(cardHolderModels))
	for _, value := range cardHolderModels {
		cardHolders = append(cardHolders, res.CardHolderResponse{
			ID:                value.ID,
			AuthorID:          value.AuthorID,
			Title:             value.Title,
			Description:       value.Description,
			Tags:              value.Tags,
			Categories:        res.ToCategoryResponse(value.Categories),
			NumberOfQuestions: len(value.Cards),
			CreatedAt:         value.CreatedAt,
		})
	}
	return cardHolders
}

func toCardHolderModel(req request.CardHolderRequest) models.CardHolder {
	card := models.CardHolder{
		AuthorID:    *req.AuthorID,
		Title:       req.Title,
		Description: *req.Description,
		Tags:        req.Tags,
		Cards:       toCardModels(req.Cards),
	}
	for _, catID := range req.Categories {
		card.Categories = append(card.Categories, models.Category{ID: catID})
	}
	return card
}

func toCardModels(cards []request.CardRequest) []models.Card {
	modelsCards := make([]models.Card, 0, len(cards))
	for _, value := range cards {
		modelsCards = append(modelsCards, models.Card{
			Question: value.Question,
			Answer:   value.Answer,
		})
	}
	return modelsCards
}

func toCardHolderDetailResponse(cardHolder *models.CardHolder) *res.CardHolderDetailResponse {
	return &res.CardHolderDetailResponse{
		ID:          cardHolder.ID,
		AuthorID:    cardHolder.AuthorID,
		Title:       cardHolder.Title,
		Description: cardHolder.Description,
		Tags:        cardHolder.Tags,
		Categories:  res.ToCategoryResponse(cardHolder.Categories),
		Cards:       toCardResponse(cardHolder.Cards),
		CreatedAt:   cardHolder.CreatedAt,
	}
}

func toCardResponse(cards []models.Card) []res.CardResponse {
	newCards := make([]res.CardResponse, 0, len(cards))
	for _, value := range cards {
		newCards = append(newCards, res.CardResponse{
			ID:       value.ID,
			Question: value.Question,
			Answer:   value.Answer,
		})
	}
	return newCards
}
