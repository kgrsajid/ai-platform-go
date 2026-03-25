package store

import (
	"project-go/internal/repository/card"
	"project-go/internal/repository/chat"
	"project-go/internal/repository/passwordreset"
	"project-go/internal/repository/progression"
	"project-go/internal/repository/session"
	"project-go/internal/repository/test"
	"project-go/internal/repository/testcategory"
	"project-go/internal/repository/user"

	"gorm.io/gorm"
)

type Store struct {
	UserRepo          *user.UserRepository
	SessionRepo       *session.SessionRepository
	ChatRepo          *chat.ChatRepository
	TestRepo          *testrepo.TestRepository
	CardRepo          *card.CardRepository
	QuestionRepo      *testrepo.QuestionRepository
	CategoryRepo      *testcategory.TestCategoryRepository
	TestViewRepo      *testrepo.TestViewRepo
	PasswordResetRepo *passwordreset.Repository
	ProgressionRepo   *progression.Repository
}

func NewStore(db *gorm.DB) *Store {
	return &Store{
		UserRepo:          user.NewUserRepo(db),
		SessionRepo:       session.NewSessionRepo(db),
		ChatRepo:          chat.NewChatRepo(db),
		TestRepo:          testrepo.NewTestRepo(db),
		CardRepo:          card.New(db),
		QuestionRepo:      testrepo.NewQuestionRepo(db),
		CategoryRepo:      testcategory.New(db),
		TestViewRepo:      testrepo.NewTestViewRepo(db),
		PasswordResetRepo: passwordreset.NewRepo(db),
		ProgressionRepo:   progression.NewRepository(db),
	}
}
