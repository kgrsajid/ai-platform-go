package models

import (
	"time"

	"github.com/lib/pq"
	_ "gorm.io/gorm"
)

type Role string

type Status string

const (
	Admin   Role = "admin"
	Teacher Role = "teacher"
	Student Role = "student"
)

const (
	Pending Status = "pending"
	Success Status = "success"
	Error   Status = "error"
)

type Difficulty string

const (
	Easy   Difficulty = "easy"
	Medium Difficulty = "medium"
	Hard   Difficulty = "hard"
)

func IsValidRole(r string) bool {
	switch Role(r) {
	case Admin, Teacher, Student:
		return true
	}
	return false
}

// Пользователь
type User struct {
	ID        uint    `gorm:"primaryKey" json:"id"`
	Name      string  `gorm:"size:100;not null" json:"name"`
	Email     string  `gorm:"size:100;unique;not null" json:"email"`
	Password  string  `gorm:"size:255;not null" json:"-"`
	Role      Role    `gorm:"type:role_enum;not null" json:"role"`                                                         // admin / teacher / student
	Grade     int     `gorm:"index;default:null" json:"grade"`                                                             // 0-11 (null for non-students)
	School    string  `gorm:"size:255;default:''" json:"school"`
	Avatar    string  `gorm:"size:500;default:''" json:"avatar"`
	Language  string  `gorm:"size:5;default:'en'" json:"language"`                                                        // 'en', 'ru', 'kz'
	Students  []*User `gorm:"many2many:teachers_students;joinForeignKey:TeacherID;JoinReferences:StudentID"` // для учителей
	Teachers  []*User `gorm:"many2many:teachers_students;joinForeignKey:StudentID;JoinReferences:TeacherID"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Связь учителей и учеников
type TeacherStudent struct {
	TeacherID uint `gorm:"primaryKey"`
	StudentID uint `gorm:"primaryKey"`
}

type SessionHistory struct {
	ID        uint `gorm:"primaryKey"`
	StudentID uint
	Student   User   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Title     string `gorm:"size:255;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Chats     []ChatMessage `gorm:"foreignKey:SessionID"`
}

type ChatMessage struct {
	ID        uint `gorm:"primaryKey"`
	SessionID uint
	Session   SessionHistory `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	// "user" или "bot"
	Role string `gorm:"type:varchar(10)"`

	Status Status `gorm:type:status_enum;not_null;default:'success'`
	// Сам текст сообщения
	Content string `gorm:"type:text"`

	CreatedAt time.Time
}

// Категория
type Category struct {
	ID          uint         `gorm:"primaryKey" json:"id"`
	Name        string       `gorm:"size:100;unique;not null" json:"name"`
	MinGrade    int          `gorm:"index;default:0" json:"min_grade"`
	MaxGrade    int          `gorm:"index;default:11" json:"max_grade"`
	Tests       []Test       `gorm:"many2many:test_categories" json:"tests"`
	CardHolders []CardHolder `gorm:"many2many:card_categories" json:"card_holders"`
}

// Тест
type Test struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Title       string         `gorm:"size:255;not null" json:"title"`
	Description string         `gorm:"type:text" json:"description"`
	Difficulty  Difficulty     `gorm:"type:difficulty_enum;not null;default:'medium'" json:"difficulty"`
	MinGrade    int            `gorm:"index;default:0" json:"min_grade"`
	MaxGrade    int            `gorm:"index;default:11" json:"max_grade"`
	Categories  []Category     `gorm:"many2many:test_categories" json:"categories"`
	Subjects    []Subject      `gorm:"many2many:test_subjects" json:"subjects"`
	Tags        pq.StringArray `gorm:"type:text[]" json:"tags"`
	Questions   []TestQuestion `gorm:"foreignKey:TestID" json:"questions"`
	IsPrivate   bool           `gorm:"not null;default:false;index" json:"is_private"`
	AuthorID    uint           `gorm:"not null;index" json:"author_id"`
	Author      User           `gorm:"foreignKey:AuthorID" json:"author"`
	ViewCount   uint           `gorm:"default:0" json:"view_count"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

type TestView struct {
	ID        uint
	TestID    uint
	UserID    uint
	CreatedAt time.Time
}

// Вопросы теста
type TestQuestion struct {
	ID          uint `gorm:"primaryKey"`
	TestID      uint
	Test        Test   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Question    string `gorm:"type:text;not null"`
	DurationSec int
	Options     []TestOption `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

// Результаты теста
type TestResult struct {
	ID        uint `gorm:"primaryKey"`
	StudentID uint
	Student   User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	TestID uint
	Test   Test `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	Score      int     // набранные баллы
	MaxScore   int     // максимум
	Percentage float64 // 0–100

	Attempt int // номер попытки

	StartedAt   time.Time
	FinishedAt  time.Time
	DurationSec int

	CreatedAt time.Time
}

type TestOption struct {
	ID             uint `gorm:"primaryKey"`
	TestQuestionID uint
	OptionText     string `gorm:"type:text;not null"`
	IsCorrect      bool   `gorm:"default:false"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// Карточки
type CardHolder struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	AuthorID    uint           `json:"author_id"`
	Author      User           `gorm:"foreignKey:AuthorID" json:"author,omitempty"`
	Title       string         `gorm:"type:text;not null" json:"title"`
	Description string         `gorm:"type:text" json:"description"`
	IsPrivate   bool           `gorm:"not null;default:false;index" json:"is_private"`
	MinGrade    int            `gorm:"index;default:0" json:"min_grade"`
	MaxGrade    int            `gorm:"index;default:11" json:"max_grade"`
	Tags        pq.StringArray `gorm:"type:text[]" json:"tags"`
	Categories  []Category     `gorm:"many2many:card_categories" json:"categories"`
	Cards       []Card         `json:"cards"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

type Card struct {
	ID           uint `gorm:"primaryKey"`
	CardHolderID uint
	CardHolder   CardHolder `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Question     string     `gorm:"type:text;not null"`
	Answer       string     `gorm:"type:text;not null"`
}

type PasswordReset struct {
	ID        uint      `gorm:"primaryKey"`
	Email     string    `gorm:"size:100;not null;index"`
	Code      string    `gorm:"size:10;not null"`
	Token     string    `gorm:"size:255;index"`
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time
}

// Игры
type Game struct {
	ID          uint   `gorm:"primaryKey"`
	Title       string `gorm:"type:text;not null"`
	Description string `gorm:"type:text"`
	CreatedAt   time.Time
}

type GameResult struct {
	ID        uint `gorm:"primaryKey"`
	StudentID uint
	Student   User `gorm:"constraint:onUpdate:CASCADE,onDelete:CASCADE"`
	GameID    uint
	Game      Game `gorm:"constraint:onUpdate:CASCADE,onDelete:CASCADE"`
	Score     int  `gorm:"not null"`
	CreatedAt time.Time
}
