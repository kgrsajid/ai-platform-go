package store

import (
	"project-go/internal/http-server/repository/chat"
	"project-go/internal/http-server/repository/session"
	testcategory "project-go/internal/http-server/repository/test-category"
	"project-go/internal/http-server/repository/test/question"
	"project-go/internal/http-server/repository/test/test"
	"project-go/internal/http-server/repository/test/view"
	"project-go/internal/http-server/repository/user"

	"gorm.io/gorm"
)

type Store struct {
	UserRepo     *user.UserRepository
	SessionRepo  *session.SessionRepository
	ChatRepo     *chat.ChatRepository
	TestRepo     *test.TestRepository
	QuestionRepo *question.QuestionRepository
	CategoryRepo *testcategory.TestCategoryRepository
	TestViewRepo *view.TestViewRepo
}

func NewStore(db *gorm.DB) *Store {
	return &Store{
		UserRepo:     user.NewUserRepo(db),
		SessionRepo:  session.NewSessionRepo(db),
		ChatRepo:     chat.NewChatnRepo(db),
		TestRepo:     test.NewTestRepo(db),
		QuestionRepo: question.NewQuestionRepo(db),
		CategoryRepo: testcategory.New(db),
		TestViewRepo: view.NewTestView(db),
	}
}
