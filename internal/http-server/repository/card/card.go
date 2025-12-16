package card

import (
	"project-go/internal/models"

	"gorm.io/gorm"
)

type CardRepository struct {
	db *gorm.DB
}

func New(db *gorm.DB) *CardRepository {
	return &CardRepository{
		db: db,
	}
}

func (r *CardRepository) CreateCard(card *models.Card) (*models.Card, error) {
	if err := r.db.Create(card).Error; err != nil {
		return nil, err
	}
	if err := r.db.
		Preload("CardHolder").
		Preload("CardHolder.Student").
		First(card, card.ID).Error; err != nil {
		return nil, err
	}

	return card, nil
}

func (r *CardRepository) GetCardsByCardHolderId(cardHolderId uint) ([]models.Card, error) {
	var cards []models.Card
	if err := r.db.
		Preload("CardHolder").
		Preload("CardHolder.Student").
		Where("cardHolder_id = ?", cardHolderId).
		Order("created_at DESC").
		Find(&cards).Error; err != nil {
		return nil, err
	}
	return cards, nil
}
