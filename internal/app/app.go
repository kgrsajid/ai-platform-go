package app

import (
	"log/slog"
	"net/http"
	chatAdd "project-go/internal/http-server/handlers/chat/add"
	addbycreatingsession "project-go/internal/http-server/handlers/chat/addByCreatingSession"
	GetChatBySessionId "project-go/internal/http-server/handlers/chat/getBySessionId"
	GetAllSessions "project-go/internal/http-server/handlers/session/getAll"
	CategoryCreate "project-go/internal/http-server/handlers/test-category/create"
	CategoryGetAll "project-go/internal/http-server/handlers/test-category/getAll"
	AddTestView "project-go/internal/http-server/handlers/test/addTestView"
	AddTestResult "project-go/internal/http-server/handlers/test/test/addResult"
	TestCreate "project-go/internal/http-server/handlers/test/test/create"
	TestGetAll "project-go/internal/http-server/handlers/test/test/getAll"
	GetAllUserTestResults "project-go/internal/http-server/handlers/test/test/getAllUserResults"
	TestGetById "project-go/internal/http-server/handlers/test/test/getById"
	UserCreate "project-go/internal/http-server/handlers/user/create"
	UserLogin "project-go/internal/http-server/handlers/user/login"
	"project-go/internal/http-server/repository/store"
	chatservice "project-go/internal/http-server/service/chat"
	sessionService "project-go/internal/http-server/service/session"
	testservice "project-go/internal/http-server/service/test"
	testcategoryservice "project-go/internal/http-server/service/test-category"
	testviewservice "project-go/internal/http-server/service/test-view"
	userservice "project-go/internal/http-server/service/user"
	"project-go/internal/lib/jwt"
)

type App struct {
	AddChatHandler              http.HandlerFunc
	AddMessageByCreatingSession http.HandlerFunc
	GetChatBySessionIdHandler   http.HandlerFunc
	GetAllSessions              http.HandlerFunc
	TestCreate                  http.HandlerFunc
	TestGetAll                  http.HandlerFunc
	TestGetById                 http.HandlerFunc
	TestResultsGetALl           http.HandlerFunc
	TestResultsAdd              http.HandlerFunc
	UserCreate                  http.HandlerFunc
	UserLogin                   http.HandlerFunc
	CreateCategory              http.HandlerFunc
	GetAllCategories            http.HandlerFunc
	TestViewAdd                 http.HandlerFunc
}

func New(log *slog.Logger, store *store.Store, jwtKey string) *App {
	// репозитории
	chatRepo := store.ChatRepo
	sessionRepo := store.SessionRepo
	testRepo := store.TestRepo
	questionRepo := store.QuestionRepo
	userRepo := store.UserRepo
	categoryRepo := store.CategoryRepo
	testViewRepo := store.TestViewRepo
	// сервисы
	chatService := chatservice.New(chatRepo, sessionRepo)
	sessionService := sessionService.New(sessionRepo)
	testService := testservice.New(testRepo, questionRepo)
	userService := userservice.New(userRepo)
	testCategoryService := testcategoryservice.New(categoryRepo)
	testViewService := testviewservice.New(testViewRepo)
	// хендлеры
	AddChatHandler := chatAdd.New(log, chatService)
	AddMessageByCreatingSession := addbycreatingsession.New(log, chatService)
	GetChatBySessionIdHandler := GetChatBySessionId.New(log, chatService)
	GetAllSessions := GetAllSessions.New(log, sessionService)
	TestCreate := TestCreate.New(log, testService)
	TestGetAll := TestGetAll.New(log, testService)
	TestGetById := TestGetById.New(log, testService)
	TestResultsGetAll := GetAllUserTestResults.New(log, testService)
	AddTestResult := AddTestResult.New(log, testService)
	AddTestView := AddTestView.New(log, testViewService)
	CategoryCreate := CategoryCreate.New(log, testCategoryService)
	CategoryGetAll := CategoryGetAll.New(log, *testCategoryService)
	UserCreate := UserCreate.New(log, userService)
	UserLogin := UserLogin.New(log, userService, jwt.NewJWTService(jwtKey))

	return &App{
		AddChatHandler:              AddChatHandler,
		AddMessageByCreatingSession: AddMessageByCreatingSession,
		GetChatBySessionIdHandler:   GetChatBySessionIdHandler,
		GetAllSessions:              GetAllSessions,
		TestCreate:                  TestCreate,
		TestGetAll:                  TestGetAll,
		TestGetById:                 TestGetById,
		TestResultsAdd:              AddTestResult,
		TestResultsGetALl:           TestResultsGetAll,
		CreateCategory:              CategoryCreate,
		GetAllCategories:            CategoryGetAll,
		UserCreate:                  UserCreate,
		UserLogin:                   UserLogin,
		TestViewAdd:                 AddTestView,
	}
}
