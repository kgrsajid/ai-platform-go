package card

import (
	"fmt"
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
	var categories []models.Category
	for _, c := range card.Categories {
		categories = append(categories, models.Category{
			ID: c.ID,
		})
	}
	if len(categories) > 0 {
		if err := r.db.Where("id IN ?", extractIDs(categories)).Find(&categories).Error; err != nil {
			return nil, err
		}
	}
	card.Categories = categories
	if err := r.db.Create(card).Error; err != nil {
		return nil, err
	}

	return card, nil
}

func extractIDs(cats []models.Category) []uint {
	ids := make([]uint, len(cats))
	for i, c := range cats {
		ids[i] = c.ID
	}
	return ids
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

func (r *CardRepository) UpdateCard(cardholder *models.CardHolder) (*models.CardHolder, error) {
	cx := r.db.Begin()

	if cx.Error != nil {
		return nil, cx.Error
	}

	fmt.Println("cardholder", cardholder.ID)

	if err := cx.Model(&models.CardHolder{}).
		Where("id = ?", cardholder.ID).
		Updates(map[string]interface{}{
			"title":       cardholder.Title,
			"description": cardholder.Description,
			"tags":        cardholder.Tags,
		}).Error; err != nil {
		cx.Rollback()
		return nil, err
	}
	var categories []models.Category
	if len(cardholder.Categories) > 0 {
		var categoryIDs []uint
		for _, c := range cardholder.Categories {
			categoryIDs = append(categoryIDs, c.ID)
		}

		if err := cx.
			Where("id IN ?", categoryIDs).
			Find(&categories).Error; err != nil {
			cx.Rollback()
			return nil, err
		}
	}

	if err := cx.
		Model(&models.CardHolder{ID: cardholder.ID}).
		Association("Categories").
		Replace(categories); err != nil {
		cx.Rollback()
		return nil, err
	}
	if err := cx.
		Where("card_holder_id = ?", cardholder.ID).
		Delete(&models.Card{}).Error; err != nil {
		cx.Rollback()
		return nil, err
	}
	for i := range cardholder.Cards {
		cardholder.Cards[i].CardHolderID = cardholder.ID
	}

	if len(cardholder.Cards) > 0 {
		if err := cx.Create(&cardholder.Cards).Error; err != nil {
			cx.Rollback()
			return nil, err
		}
	}
	// 5️⃣ Commit
	if err := cx.Commit().Error; err != nil {
		return nil, err
	}

	return cardholder, nil

}

func (r *CardRepository) GetCardsByCardHolderId(cardHolderId uint) (*models.CardHolder, error) {
	var cards models.CardHolder
	if err := r.db.Where("id = ?", cardHolderId).
		Preload("Categories").
		Preload("Cards").
		Order("created_at DESC").
		Find(&cards).Error; err != nil {
		return nil, err
	}
	return &cards, nil
}
