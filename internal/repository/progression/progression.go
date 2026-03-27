package progression

import (
	"errors"
	"fmt"
	"time"

	"project-go/internal/models"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// GetOrCreateProgress returns the user's progress, creating it if it doesn't exist
func (r *Repository) GetOrCreateProgress(userID uint) (*models.UserProgress, *models.User, error) {
	var progress models.UserProgress
	err := r.db.Preload("User").Where("user_id = ?", userID).First(&progress).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			progress = models.UserProgress{
				UserID: userID,
			}
			// Also get user for grade
			var user models.User
			r.db.First(&user, userID)
			progress.User = user
			if err := r.db.Create(&progress).Error; err != nil {
				return nil, nil, fmt.Errorf("create progress: %w", err)
			}
			return &progress, &user, nil
		}
		return nil, nil, fmt.Errorf("get progress: %w", err)
	}
	return &progress, &progress.User, nil
}

// AddPoints adds points and creates a transaction record
func (r *Repository) AddPoints(userID uint, amount int, source, referenceID, description string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Update progress
		result := tx.Model(&models.UserProgress{}).
			Where("user_id = ?", userID).
			Updates(map[string]interface{}{
				"total_points":     gorm.Expr("total_points + ?", amount),
				"available_points": gorm.Expr("available_points + ?", amount),
			})
		if result.Error != nil {
			return result.Error
		}

		// If no rows affected, create progress first
		if result.RowsAffected == 0 {
			var existing models.UserProgress
			if err := tx.Where("user_id = ?", userID).First(&existing).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					progress := models.UserProgress{UserID: userID, TotalPoints: amount, AvailablePoints: amount}
					if err := tx.Create(&progress).Error; err != nil {
						return err
					}
				} else {
					return err
				}
			} else {
				if err := tx.Model(&existing).Updates(map[string]interface{}{
					"total_points":     gorm.Expr("total_points + ?", amount),
					"available_points": gorm.Expr("available_points + ?", amount),
				}).Error; err != nil {
					return err
				}
			}
		}

		// Create transaction log
		transaction := models.PointTransaction{
			UserID:      userID,
			Amount:      amount,
			Source:      source,
			ReferenceID: referenceID,
			Description: description,
		}
		return tx.Create(&transaction).Error
	})
}

// SpendPoints deducts points and creates a transaction record
func (r *Repository) SpendPoints(userID uint, amount int, source, referenceID, description string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Check balance
		var progress models.UserProgress
		if err := tx.Where("user_id = ?", userID).First(&progress).Error; err != nil {
			return fmt.Errorf("progress not found: %w", err)
		}

		if progress.AvailablePoints < amount {
			return fmt.Errorf("insufficient points: have %d, need %d", progress.AvailablePoints, amount)
		}

		// Deduct points
		if err := tx.Model(&progress).Updates(map[string]interface{}{
			"available_points": gorm.Expr("available_points - ?", amount),
		}).Error; err != nil {
			return err
		}

		// Create transaction log (negative amount)
		transaction := models.PointTransaction{
			UserID:      userID,
			Amount:      -amount,
			Source:      source,
			ReferenceID: referenceID,
			Description: description,
		}
		return tx.Create(&transaction).Error
	})
}

// AddXP adds XP and handles level-up
func (r *Repository) AddXP(userID uint, xpAmount int, source string) (*models.UserProgress, error) {
	var progress models.UserProgress
	err := r.db.Transaction(func(tx *gorm.DB) error {
		// Get or create progress with user preloaded
		if err := tx.Preload("User").Where("user_id = ?", userID).First(&progress).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				progress = models.UserProgress{UserID: userID}
				if err := tx.Create(&progress).Error; err != nil {
					return err
				}
				// Load user for grade
				tx.Preload("User").First(&progress, progress.ID)
			} else {
				return err
			}
		}

		// Determine grade band for level calculation
		grade := progress.User.Grade
		if grade <= 0 || grade > 11 {
			grade = 5 // default to Explorers
		}
		band := models.GetGradeBand(grade)
		xpPerLevel := band.XPPerLevel()

		// Add XP
		newXP := progress.TotalXP + xpAmount
		newLevel := progress.CurrentLevel

		// Check for level-ups
		xpForCurrentLevel := 0
		for i := 1; i < newLevel; i++ {
			xpForCurrentLevel += xpPerLevel
		}
		remainingXP := newXP - xpForCurrentLevel

		for remainingXP >= xpPerLevel && newLevel < band.MaxLevel() {
			newLevel++
			remainingXP -= xpPerLevel
		}

		// Update progress
		progress.TotalXP = newXP
		progress.CurrentLevel = newLevel

		if err := tx.Save(&progress).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	return &progress, nil
}

// RecordActivity handles daily activity tracking and streak calculation
func (r *Repository) RecordActivity(userID uint) (*models.DailyActivity, *models.UserProgress, error) {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	var dailyActivity models.DailyActivity
	var progress models.UserProgress

	err := r.db.Transaction(func(tx *gorm.DB) error {
		// Get or create progress
		if err := tx.Where("user_id = ?", userID).First(&progress).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				progress = models.UserProgress{UserID: userID}
				if err := tx.Create(&progress).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		}

		// Get or create today's activity
		err := tx.Where("user_id = ? AND date = ?", userID, today).First(&dailyActivity).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				dailyActivity = models.DailyActivity{
					UserID: userID,
					Date:   today,
				}
				if err := tx.Create(&dailyActivity).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		}

		// Calculate streak
		newStreak := 1
		if progress.LastActiveAt != nil {
			lastDay := time.Date(progress.LastActiveAt.Year(), progress.LastActiveAt.Month(), progress.LastActiveAt.Day(), 0, 0, 0, 0, progress.LastActiveAt.Location())
			diff := today.Sub(lastDay).Hours() / 24

			switch {
			case diff == 0:
				// Same day, no change
				newStreak = progress.CurrentStreak
			case diff == 1:
				// Consecutive day
				newStreak = progress.CurrentStreak + 1
			default:
				// Gap — reset streak
				newStreak = 1
			}
		}

		// Update streak
		if newStreak > progress.LongestStreak {
			progress.LongestStreak = newStreak
		}
		progress.CurrentStreak = newStreak
		nowPtr := now
		progress.LastActiveAt = &nowPtr

		if err := tx.Save(&progress).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, nil, err
	}

	return &dailyActivity, &progress, nil
}

// GetTransactions returns paginated point transactions
func (r *Repository) GetTransactions(userID uint, page, limit int) ([]models.PointTransaction, int64, error) {
	var transactions []models.PointTransaction
	var total int64

	offset := (page - 1) * limit

	db := r.db.Where("user_id = ?", userID)
	db.Model(&models.PointTransaction{}).Count(&total)

	err := db.Order("created_at DESC").Offset(offset).Limit(limit).Find(&transactions).Error
	return transactions, total, err
}

// GetRewards returns active rewards filtered by grade
func (r *Repository) GetRewards(grade int) ([]models.Reward, error) {
	var rewards []models.Reward
	err := r.db.Where("is_active = ? AND min_grade <= ? AND max_grade >= ?", true, grade, grade).
		Order("point_cost ASC").
		Find(&rewards).Error
	return rewards, err
}

// RedeemReward redeems a reward for a user
func (r *Repository) RedeemReward(userID uint, rewardID uint) (*models.Redemption, error) {
	var redemption models.Redemption

	err := r.db.Transaction(func(tx *gorm.DB) error {
		// Get reward
		var reward models.Reward
		if err := tx.Where("id = ? AND is_active = ?", rewardID, true).First(&reward).Error; err != nil {
			return fmt.Errorf("reward not found: %w", err)
		}

		// Check stock
		if reward.TotalStock != -1 && reward.TotalStock <= 0 {
			return fmt.Errorf("reward out of stock")
		}

		// Check user balance
		var progress models.UserProgress
		if err := tx.Where("user_id = ?", userID).First(&progress).Error; err != nil {
			return fmt.Errorf("progress not found: %w", err)
		}

		if progress.AvailablePoints < reward.PointCost {
			return fmt.Errorf("insufficient points: have %d, need %d", progress.AvailablePoints, reward.PointCost)
		}

		// Generate coupon code
		couponCode := fmt.Sprintf("SDU-%d-%d-%s", userID, rewardID, time.Now().Format("20060102"))

		// Create redemption
		redemption = models.Redemption{
			UserID:      userID,
			RewardID:    rewardID,
			PointsSpent: reward.PointCost,
			CouponCode:  couponCode,
		}

		if err := tx.Create(&redemption).Error; err != nil {
			return err
		}

		// Deduct points
		if err := tx.Model(&progress).Update("available_points", gorm.Expr("available_points - ?", reward.PointCost)).Error; err != nil {
			return err
		}

		// Log transaction
		transaction := models.PointTransaction{
			UserID:      userID,
			Amount:      -reward.PointCost,
			Source:      "redemption",
			ReferenceID: fmt.Sprintf("%d", rewardID),
			Description: fmt.Sprintf("Redeemed: %s", reward.Title),
		}
		if err := tx.Create(&transaction).Error; err != nil {
			return err
		}

		// Decrease stock
		if reward.TotalStock != -1 {
			tx.Model(&reward).Update("total_stock", gorm.Expr("total_stock - 1"))
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Preload reward
	r.db.Preload("Reward").First(&redemption, redemption.ID)

	return &redemption, nil
}

// GetMyRedemptions returns user's redeemed coupons
func (r *Repository) GetMyRedemptions(userID uint) ([]models.Redemption, error) {
	var redemptions []models.Redemption
	err := r.db.Where("user_id = ?", userID).Preload("Reward").Order("created_at DESC").Find(&redemptions).Error
	return redemptions, err
}

// GetSubjects returns all subjects
func (r *Repository) GetSubjects() ([]models.Subject, error) {
	var subjects []models.Subject
	err := r.db.Order("sort_order ASC").Find(&subjects).Error
	return subjects, err
}
