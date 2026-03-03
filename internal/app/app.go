package app

import (
	"log/slog"
	"net/http"
	client "project-go/internal/http-server/client/chat"
	CardCreate "project-go/internal/http-server/handlers/card/create"
	CardGetAll "project-go/internal/http-server/handlers/card/getAll"
	CardGetById "project-go/internal/http-server/handlers/card/getById"
	CardUpdate "project-go/internal/http-server/handlers/card/update"
	chatAdd "project-go/internal/http-server/handlers/chat/add"
	addbycreatingsession "project-go/internal/http-server/handlers/chat/addByCreatingSession"
	GetChatBySessionId "project-go/internal/http-server/handlers/chat/getBySessionId"
	RetryLastMessage "project-go/internal/http-server/handlers/chat/retry"
	CreateSession "project-go/internal/http-server/handlers/session/create"
	DeleteSession "project-go/internal/http-server/handlers/session/delete"
	GetAllSessions "project-go/internal/http-server/handlers/session/getAll"
	CategoryCreate "project-go/internal/http-server/handlers/test-category/create"
	CategoryGetAll "project-go/internal/http-server/handlers/test-category/getAll"
	AddTestView "project-go/internal/http-server/handlers/test/addTestView"
	AddTestResult "project-go/internal/http-server/handlers/test/test/addResult"
	TestCreate "project-go/internal/http-server/handlers/test/test/create"
	TestGetAll "project-go/internal/http-server/handlers/test/test/getAll"
	GetAllUserTestResults "project-go/internal/http-server/handlers/test/test/getAllUserResults"
	TestGetById "project-go/internal/http-server/handlers/test/test/getById"
	TestUpdate "project-go/internal/http-server/handlers/test/test/updateTest"
	UserCreate "project-go/internal/http-server/handlers/user/create"
	UserLogin "project-go/internal/http-server/handlers/user/login"
	"project-go/internal/http-server/repository/store"
	cardservice "project-go/internal/http-server/service/card"
	chatservice "project-go/internal/http-server/service/chat"
	sessionService "project-go/internal/http-server/service/session"
	testservice "project-go/internal/http-server/service/test"
	testcategoryservice "project-go/internal/http-server/service/test-category"
	testviewservice "project-go/internal/http-server/service/test-view"
	userservice "project-go/internal/http-server/service/user"
	websocket "project-go/internal/http-server/websocket/chat"
	"project-go/internal/lib/auth"
	"project-go/internal/lib/jwt"
)

type App struct {
	AddChatHandler              http.HandlerFunc
	AddMessageByCreatingSession http.HandlerFunc
	GetChatBySessionIdHandler   http.HandlerFunc
	GetAllSessions              http.HandlerFunc
	CreateSession               http.HandlerFunc
	DeleteSession               http.HandlerFunc
	RetryLastMessage            http.HandlerFunc
	TestCreate                  http.HandlerFunc
	TestUpdate                  http.HandlerFunc
	TestGetAll                  http.HandlerFunc
	TestGetById                 http.HandlerFunc
	TestResultsGetALl           http.HandlerFunc
	TestResultsAdd              http.HandlerFunc
	UserCreate                  http.HandlerFunc
	UserLogin                   http.HandlerFunc
	CreateCategory              http.HandlerFunc
	GetAllCategories            http.HandlerFunc
	TestViewAdd                 http.HandlerFunc
	CardCreate                  http.HandlerFunc
	CardGetAll                  http.HandlerFunc
	CardGetById                 http.HandlerFunc
	CardUpdate                  http.HandlerFunc

	WSAddMessage *websocket.Handler
}

func New(log *slog.Logger, store *store.Store, jwtKey string, aiUrl string) *App {
	// репозитории
	chatRepo := store.ChatRepo
	sessionRepo := store.SessionRepo
	testRepo := store.TestRepo
	questionRepo := store.QuestionRepo
	userRepo := store.UserRepo
	categoryRepo := store.CategoryRepo
	testViewRepo := store.TestViewRepo
	cardRepo := store.CardRepo

	//client
	aiChat := client.NewAIClient(aiUrl)

	//web socket
	hub := websocket.NewHub()
	authService := auth.New(jwtKey)
	// сервисы
	chatService := chatservice.New(chatRepo, sessionRepo, aiChat)
	sessionService := sessionService.New(sessionRepo)
	testService := testservice.New(testRepo, questionRepo)
	cardService := cardservice.New(cardRepo)
	userService := userservice.New(userRepo)
	testCategoryService := testcategoryservice.New(categoryRepo)
	testViewService := testviewservice.New(testViewRepo)
	// хендлеры
	AddChatHandler := chatAdd.New(log, chatService)
	AddMessageByCreatingSession := addbycreatingsession.New(log, chatService)
	GetChatBySessionIdHandler := GetChatBySessionId.New(log, chatService)
	GetAllSessions := GetAllSessions.New(log, sessionService)
	CreateSession := CreateSession.New(log, sessionService)
	DeleteSession := DeleteSession.New(log, sessionService)
	RetryLastMessage := RetryLastMessage.New(log, chatService)
	TestCreate := TestCreate.New(log, testService)
	TestUpdate := TestUpdate.New(log, testService)
	TestGetAll := TestGetAll.New(log, testService)
	TestGetById := TestGetById.New(log, testService)
	TestResultsGetAll := GetAllUserTestResults.New(log, testService)
	AddTestResult := AddTestResult.New(log, testService)
	AddTestView := AddTestView.New(log, testViewService)
	CategoryCreate := CategoryCreate.New(log, testCategoryService)
	CategoryGetAll := CategoryGetAll.New(log, *testCategoryService)
	CardCreate := CardCreate.New(log, cardService)
	CardGetAll := CardGetAll.New(log, cardService)
	CardGetById := CardGetById.New(log, cardService)
	CardUpdate := CardUpdate.New(log, cardService)
	UserCreate := UserCreate.New(log, userService)
	UserLogin := UserLogin.New(log, userService, jwt.NewJWTService(jwtKey))

	wsHandler := websocket.NewHandler(hub, authService, chatService)

	return &App{
		AddChatHandler:              AddChatHandler,
		AddMessageByCreatingSession: AddMessageByCreatingSession,
		GetChatBySessionIdHandler:   GetChatBySessionIdHandler,
		GetAllSessions:              GetAllSessions,
		CreateSession:               CreateSession,
		DeleteSession:               DeleteSession,
		RetryLastMessage:            RetryLastMessage,
		TestCreate:                  TestCreate,
		TestUpdate:                  TestUpdate,
		TestGetAll:                  TestGetAll,
		TestGetById:                 TestGetById,
		TestResultsAdd:              AddTestResult,
		TestResultsGetALl:           TestResultsGetAll,
		CreateCategory:              CategoryCreate,
		GetAllCategories:            CategoryGetAll,
		CardCreate:                  CardCreate,
		CardGetAll:                  CardGetAll,
		CardGetById:                 CardGetById,
		CardUpdate:                  CardUpdate,
		UserCreate:                  UserCreate,
		UserLogin:                   UserLogin,
		TestViewAdd:                 AddTestView,
		WSAddMessage:                wsHandler,
	}
}
