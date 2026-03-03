package testrepo

import (
	"project-go/internal/models"

	"gorm.io/gorm"
)

type TestViewRepo struct {
	db *gorm.DB
}

func NewTestViewRepo(db *gorm.DB) *TestViewRepo {
	return &TestViewRepo{db: db}
}

func (r *TestViewRepo) AddTestView(testId uint, userId uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var count int64
		if err := tx.Model(&models.TestView{}).
			Where("test_id = ? AND user_id = ?", testId, userId).
			Count(&count).Error; err != nil {
			return err
		}

		if count > 0 {
			return nil
		}

		view := models.TestView{
			TestID: testId,
			UserID: userId,
		}
		if err := tx.Create(&view).Error; err != nil {
			return err
		}

		if err := tx.Model(&models.Test{}).
			Where("id = ?", testId).
			UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).Error; err != nil {
			return err
		}

		return nil
	})
}
