package trainer

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

// GetTrainerProfile returns the robot profile with computed stats
func (r *Repository) GetTrainerProfile(userID uint) (*models.UserProgress, *models.User, error) {
	var progress models.UserProgress
	err := r.db.Preload("User").Where("user_id = ?", userID).First(&progress).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Create default progress
			var user models.User
			r.db.First(&user, userID)
			progress = models.UserProgress{
				UserID:     userID,
				RobotName:  getDefaultRobotName(user.Grade),
				RobotColor: "#6366f1",
			}
			if err := r.db.Create(&progress).Error; err != nil {
				return nil, nil, fmt.Errorf("create progress: %w", err)
			}
			return &progress, &user, nil
		}
		return nil, nil, fmt.Errorf("get progress: %w", err)
	}
	return &progress, &progress.User, nil
}

// UpdateTrainerProfile updates robot name and color
func (r *Repository) UpdateTrainerProfile(userID uint, name, color string) (*models.UserProgress, error) {
	var progress models.UserProgress
	err := r.db.Where("user_id = ?", userID).First(&progress).Error
	if err != nil {
		return nil, fmt.Errorf("progress not found: %w", err)
	}

	updates := map[string]interface{}{}
	if name != "" {
		updates["robot_name"] = name
	}
	if color != "" {
		updates["robot_color"] = color
	}

	if len(updates) > 0 {
		if err := r.db.Model(&progress).Updates(updates).Error; err != nil {
			return nil, fmt.Errorf("update profile: %w", err)
		}
	}

	// Reload
	r.db.First(&progress, progress.ID)
	return &progress, nil
}

// GetTrainerTimeline returns level-up history from PointTransactions
func (r *Repository) GetTrainerTimeline(userID uint, limit int) ([]models.PointTransaction, error) {
	var transactions []models.PointTransaction
	err := r.db.Where("user_id = ? AND source = ?", userID, "level_up").
		Order("created_at DESC").
		Limit(limit).
		Find(&transactions).Error
	return transactions, err
}

// GetTrainerStats returns computed stats for the robot profile
func (r *Repository) GetTrainerStats(userID uint) (map[string]interface{}, error) {
	var progress models.UserProgress
	if err := r.db.Preload("User").Where("user_id = ?", userID).First(&progress).Error; err != nil {
		return nil, fmt.Errorf("progress not found: %w", err)
	}

	// Count quizzes completed
	var quizzesCompleted int64
	r.db.Model(&models.PointTransaction{}).
		Where("user_id = ? AND source = ?", userID, "quiz").
		Count(&quizzesCompleted)

	// Count flashcards studied
	var flashcardsStudied int64
	r.db.Model(&models.PointTransaction{}).
		Where("user_id = ? AND source = ?", userID, "flashcard").
		Count(&flashcardsStudied)

	// Count assignments completed
	var assignmentsCompleted int64
	r.db.Model(&models.AssignmentSubmission{}).
		Where("user_id = ? AND is_evaluated = ?", userID, true).
		Count(&assignmentsCompleted)

	// Determine stage name based on level
	stageName := getStageName(progress.CurrentLevel)

	// Calculate XP progress to next level
	grade := progress.User.Grade
	if grade <= 0 || grade > 11 {
		grade = 5
	}
	band := models.GetGradeBand(grade)
	xpPerLevel := band.XPPerLevel()
	xpForCurrentLevel := 0
	for i := 1; i < progress.CurrentLevel; i++ {
		xpForCurrentLevel += xpPerLevel
	}
	currentLevelXP := progress.TotalXP - xpForCurrentLevel
	nextLevelXP := xpPerLevel
	progressPercent := float64(currentLevelXP) / float64(nextLevelXP) * 100
	if progressPercent < 0 {
		progressPercent = 0
	}
	if progressPercent > 100 {
		progressPercent = 100
	}

	stats := map[string]interface{}{
		"robot_name":          progress.RobotName,
		"robot_level":         progress.RobotLevel,
		"robot_color":         progress.RobotColor,
		"current_level":       progress.CurrentLevel,
		"stage_name":          stageName,
		"total_xp":            progress.TotalXP,
		"current_level_xp":    currentLevelXP,
		"next_level_xp":       nextLevelXP,
		"progress_percent":    float64(currentLevelXP) / float64(nextLevelXP) * 100,
		"quizzes_completed":   quizzesCompleted,
		"flashcards_studied":  flashcardsStudied,
		"assignments_completed": assignmentsCompleted,
		"current_streak":      progress.CurrentStreak,
		"longest_streak":      progress.LongestStreak,
		"joined_at":           progress.CreatedAt,
	}

	return stats, nil
}

// GetLeaderboard returns ranked users
func (r *Repository) GetLeaderboard(sort string, page, limit int, gradeFilter *int) ([]map[string]interface{}, int64, error) {
	var total int64

	// Build base query with subquery approach to avoid GORM table alias issues
	baseQuery := r.db.Table("user_progresses").
		Joins("JOIN users ON users.id = user_progresses.user_id")

	// Apply grade filter if specified
	if gradeFilter != nil {
		baseQuery = baseQuery.Where("users.grade = ?", *gradeFilter)
	}

	// Count total
	baseQuery.Count(&total)

	// Determine sort column
	var orderClause string
	switch sort {
	case "points":
		orderClause = "user_progresses.total_points DESC"
	case "streak":
		orderClause = "user_progresses.longest_streak DESC"
	default: // level
		orderClause = "user_progresses.current_level DESC, user_progresses.total_xp DESC"
	}

	// Pagination
	offset := (page - 1) * limit

	// Use raw query for leaderboard to avoid GORM join issues
	type leaderboardRow struct {
		UserID        uint      `json:"user_id"`
		Name          string    `json:"name"`
		Email         string    `json:"email"`
		Grade         int       `json:"grade"`
		CurrentLevel  int       `json:"current_level"`
		TotalXP       int       `json:"total_xp"`
		TotalPoints   int       `json:"total_points"`
		LongestStreak int       `json:"longest_streak"`
		RobotName     string    `json:"robot_name"`
		RobotColor    string    `json:"robot_color"`
		CreatedAt     time.Time `json:"created_at"`
	}

	var rows []leaderboardRow
	err := r.db.Table("user_progresses").
		Select("user_progresses.user_id, users.name, users.email, users.grade, user_progresses.current_level, user_progresses.total_xp, user_progresses.total_points, user_progresses.longest_streak, user_progresses.robot_name, user_progresses.robot_color, user_progresses.created_at").
		Joins("JOIN users ON users.id = user_progresses.user_id").
		Where("users.role IN ?", []string{"student", "teacher", "admin"}).
		Order(orderClause).
		Offset(offset).Limit(limit).
		Scan(&rows).Error

	if err != nil {
		return nil, 0, fmt.Errorf("get leaderboard: %w", err)
	}

	// Build results
	var results []map[string]interface{}
	for i, row := range rows {
		results = append(results, map[string]interface{}{
			"rank":           offset + i + 1,
			"user_id":        row.UserID,
			"name":           row.Name,
			"email":          row.Email,
			"grade":          row.Grade,
			"current_level":  row.CurrentLevel,
			"total_xp":       row.TotalXP,
			"total_points":   row.TotalPoints,
			"longest_streak": row.LongestStreak,
			"robot_name":     row.RobotName,
			"robot_color":    row.RobotColor,
			"created_at":     row.CreatedAt,
		})
	}

	return results, total, nil
}

func getDefaultRobotName(grade int) string {
	band := models.GetGradeBand(grade)
	switch band {
	case models.BandSprouts:
		return "RoboTutor"
	case models.BandExplorers:
		return "Navigator"
	case models.BandChampions:
		return "Mentor"
	default:
		return "AI Buddy"
	}
}

func getStageName(level int) string {
	switch {
	case level >= 20:
		return "AI Master"
	case level >= 15:
		return "Scientist"
	case level >= 10:
		return "Problem Solver"
	case level >= 5:
		return "Thinker"
	default:
		return "Beginner"
	}
}
