package card

import (
	req "project-go/internal/http-server/dto/request"
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

func (r *CardRepository) CreateCardHolder(card *models.CardHolder) (*models.CardHolder, error) {
	var categories []uint
	var resCategory []models.Category
	for _, c := range card.Categories {
		categories = append(categories, c.ID)
	}
	if len(categories) > 0 {
		if err := r.db.Where("id IN ?", categories).Find(&resCategory).Error; err != nil {
			return nil, err
		}
	}
	card.Categories = resCategory

	if err := r.db.Create(card).Error; err != nil {
		return nil, err
	}

	return card, nil
}

func (r *CardRepository) GetAll(filter *req.CardFilter) ([]models.CardHolder, int64, error) {
	var cardHolderModels []models.CardHolder
	query := r.db.Model(&models.CardHolder{}).
		Preload("Categories").
		Preload("Cards")
	if filter.Search != nil && *filter.Search != "" {
		like := "%" + *filter.Search + "%"
		query = query.Where(`
			card_holders.title ILIKE ?
			OR EXISTS (
				SELECT 1
				FROM unnest(tags) AS tag
				WHERE tag ILIKE ?
			)
		`, like, like)
	}
	if len(filter.Categories) > 0 {
		query = query.Joins("JOIN card_categories cc ON cc.card_holder_id = card_holders.id").
			Where("cc.category_id IN ?", filter.Categories)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Order("card_holders.created_at DESC").
		Limit(filter.Limit).
		Offset(filter.Offset).
		Find(&cardHolderModels).Error; err != nil {
		return nil, 0, err
	}

	return cardHolderModels, total, nil

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
