package cardservice

import (
	req "project-go/internal/http-server/dto/request"
	res "project-go/internal/http-server/dto/response"
	"project-go/internal/models"
)

type CardRepository interface {
	CreateCardHolder(card *models.CardHolder) (*models.CardHolder, error)
	GetAll(filter *req.CardFilter) ([]models.CardHolder, int64, error)
	GetCardsByCardHolderId(cardHolderId uint) (*models.CardHolder, error)
	UpdateCard(cardholder *models.CardHolder) (*models.CardHolder, error)
}
type Service struct {
	CardRepository CardRepository
}

func New(CardRepository CardRepository) *Service {
	return &Service{
		CardRepository: CardRepository,
	}
}

func (s *Service) CreateCardHolder(card req.CardHolderRequest) (*res.CardHolderDetailResponse, error) {
	modelsCard := toCardHolderModels(card)
	createdCardHolder, err := s.CardRepository.CreateCardHolder(&modelsCard)
	if err != nil {
		return nil, err
	}
	return toCardHolderDetailResponse(createdCardHolder), nil
}

func (s *Service) UpdateCardHolder(card req.CardHolderRequest, cardId uint) (*res.CardHolderDetailResponse, error) {
	modelsCard := toCardHolderModels(card)
	modelsCard.ID = cardId
	updatedCardHolder, err := s.CardRepository.UpdateCard(&modelsCard)
	if err != nil {
		return nil, err
	}
	return toCardHolderDetailResponse(updatedCardHolder), nil
}

func (s *Service) GetAll(cardFilter req.CardFilter) ([]res.CardHolderResponse, int64, error) {
	cardHolderModels, total, err := s.CardRepository.GetAll(&cardFilter)
	if err != nil {
		return nil, 0, err
	}
	return toCardHolderResponse(cardHolderModels), total, nil
}

func (s *Service) GetCardsByCardHolderId(id uint) (*res.CardHolderDetailResponse, error) {
	cardHolder, err := s.CardRepository.GetCardsByCardHolderId(id)
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

func toCardHolderModels(cardHolder req.CardHolderRequest) models.CardHolder {
	card := models.CardHolder{
		AuthorID:    *cardHolder.AuthorID,
		Title:       cardHolder.Title,
		Description: *cardHolder.Description,
		Tags:        cardHolder.Tags,
		Cards:       toCardModels(cardHolder.Cards),
	}
	for _, catID := range cardHolder.Categories {
		card.Categories = append(card.Categories, models.Category{
			ID: catID,
		})
	}
	return card
}

func toCardModels(cards []req.CardRequest) []models.Card {
	var modelsCard = make([]models.Card, 0, len(cards))
	for _, value := range cards {
		modelsCard = append(modelsCard, models.Card{
			Question: value.Question,
			Answer:   value.Answer,
		})
	}
	return modelsCard
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
	var newCards = make([]res.CardResponse, 0, len(cards))

	for _, value := range cards {
		newCards = append(newCards, res.CardResponse{
			ID:       value.ID,
			Question: value.Question,
			Answer:   value.Answer,
		})
	}
	return newCards
}
