package models

import "time"

// UserProgress is the central gamification profile for each student
type UserProgress struct {
	ID              uint       `gorm:"primaryKey" json:"id"`
	UserID          uint       `gorm:"uniqueIndex;not null" json:"user_id"`
	User            User       `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"user,omitempty"`

	// Points & Currency (cosmetic name differs per grade band: Stars/Gems/Coins)
	TotalPoints     int        `gorm:"default:0" json:"total_points"`
	AvailablePoints int        `gorm:"default:0" json:"available_points"`

	// Level & XP
	CurrentLevel    int        `gorm:"default:1" json:"current_level"`
	TotalXP         int        `gorm:"default:0" json:"total_xp"`

	// Streaks
	CurrentStreak   int        `gorm:"default:0" json:"current_streak"`
	LongestStreak   int        `gorm:"default:0" json:"longest_streak"`
	LastActiveAt    *time.Time `gorm:"index" json:"last_active_at"`

	// AI Trainer (Phase 1)
	RobotName       string     `gorm:"size:100;default:'AI Buddy'" json:"robot_name"`
	RobotLevel      int        `gorm:"default:1" json:"robot_level"`
	RobotXP         int        `gorm:"default:0" json:"robot_xp"`
	RobotColor      string     `gorm:"size:20;default:'#6366f1'" json:"robot_color"`

	// City Builder (Phase 3)
	CityStage       int        `gorm:"default:0" json:"city_stage"`
	CityResources   int        `gorm:"default:0" json:"city_resources"`

	// Adventure Map (Phase 2)
	CurrentZone     int        `gorm:"default:1" json:"current_zone"`

	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// PointTransaction is an immutable audit log for all point changes
type PointTransaction struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"index;not null" json:"user_id"`
	User        User      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"user,omitempty"`
	Amount      int       `gorm:"not null" json:"amount"` // positive=earned, negative=spent
	Source      string    `gorm:"size:50;not null;index" json:"source"` // "quiz", "flashcard", "streak_bonus", "redemption", etc.
	ReferenceID string    `gorm:"size:100" json:"reference_id"` // e.g., test result ID
	Description string    `gorm:"type:text" json:"description"`
	CreatedAt   time.Time `gorm:"index" json:"created_at"`
}

// DailyActivity tracks daily login/activity for streaks
type DailyActivity struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	UserID            uint      `gorm:"uniqueIndex:idx_user_date;not null" json:"user_id"`
	User              User      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"user,omitempty"`
	Date              time.Time `gorm:"uniqueIndex:idx_user_date;not null" json:"date"`
	QuizzesTaken      int       `gorm:"default:0" json:"quizzes_taken"`
	FlashcardsStudied int       `gorm:"default:0" json:"flashcards_studied"`
	ChatMessages      int       `gorm:"default:0" json:"chat_messages"`
	GamesPlayed       int       `gorm:"default:0" json:"games_played"`
	PointsEarned      int       `gorm:"default:0" json:"points_earned"`
	XPEarned          int       `gorm:"default:0" json:"xp_earned"`
	CreatedAt         time.Time `json:"created_at"`
}

// Reward is a coupon/reward definition (Phase 5)
type Reward struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	Title       string     `gorm:"size:255;not null" json:"title"`
	TitleKZ     string     `gorm:"size:255" json:"title_kz"`
	TitleRU     string     `gorm:"size:255" json:"title_ru"`
	Description string     `gorm:"type:text" json:"description"`
	PartnerName string     `gorm:"size:100" json:"partner_name"`
	Category    string     `gorm:"size:50;index" json:"category"` // "food", "retail", "entertainment", "education", "virtual"
	PointCost   int        `gorm:"not null" json:"point_cost"`
	ImageURL    string     `gorm:"size:500" json:"image_url"`
	IsActive    bool       `gorm:"default:true;index" json:"is_active"`
	TotalStock  int        `gorm:"default:-1" json:"total_stock"` // -1 = unlimited
	MinGrade    int        `gorm:"default:0" json:"min_grade"`
	MaxGrade    int        `gorm:"default:11" json:"max_grade"`
	CreatedAt   time.Time  `json:"created_at"`
	ExpiresAt   *time.Time `json:"expires_at"`
}

// Redemption tracks when a user redeems a reward
type Redemption struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	UserID      uint       `gorm:"index;not null" json:"user_id"`
	User        User       `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"user,omitempty"`
	RewardID    uint       `gorm:"not null;index" json:"reward_id"`
	Reward      Reward     `gorm:"foreignKey:RewardID" json:"reward"`
	PointsSpent int        `gorm:"not null" json:"points_spent"`
	CouponCode  string     `gorm:"size:50;uniqueIndex;not null" json:"coupon_code"`
	IsUsed      bool       `gorm:"default:false" json:"is_used"`
	UsedAt      *time.Time `json:"used_at"`
	CreatedAt   time.Time  `json:"created_at"`
}
