package app

import (
	"log/slog"
	"net/http"

	client "project-go/internal/client/chat"
	"project-go/internal/config"
	cardhandler "project-go/internal/handler/card"
	chathandler "project-go/internal/handler/chat"
	sessionhandler "project-go/internal/handler/session"
	testhandler "project-go/internal/handler/test"
	testcategoryhandler "project-go/internal/handler/testcategory"
	userhandler "project-go/internal/handler/user"
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
	"project-go/internal/lib/auth"
	emaillib "project-go/internal/lib/email"
	"project-go/internal/lib/jwt"
	"project-go/internal/repository/store"
	cardservice "project-go/internal/service/card"
	chatservice "project-go/internal/service/chat"
	sessionservice "project-go/internal/service/session"
	testservice "project-go/internal/service/test"
	testcategoryservice "project-go/internal/service/testcategory"
	testviewservice "project-go/internal/service/testview"
	userservice "project-go/internal/service/user"
	wschat "project-go/internal/websocket/chat"
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
	TestResultsGetAll           http.HandlerFunc
	TestResultsAdd              http.HandlerFunc
	UserCreate                  http.HandlerFunc
	UserLogin                   http.HandlerFunc
	UserForgotPassword          http.HandlerFunc
	UserVerifyCode              http.HandlerFunc
	UserResetPassword           http.HandlerFunc
	CreateCategory              http.HandlerFunc
	GetAllCategories            http.HandlerFunc
	TestViewAdd                 http.HandlerFunc
	CardCreate                  http.HandlerFunc
	CardGetAll                  http.HandlerFunc
	CardGetById                 http.HandlerFunc
	CardUpdate                  http.HandlerFunc

	WSAddMessage *wschat.Handler
}

func New(log *slog.Logger, s *store.Store, jwtKey string, aiUrl string, emailCfg config.EmailConfig) *App {
	aiClient := client.NewAIClient(aiUrl)
	authService := auth.New(jwtKey)
	jwtService := jwt.NewJWTService(jwtKey)

	hub := wschat.NewHub()

	emailSender := emaillib.New(emaillib.Config{
		SMTPHost: emailCfg.SMTPHost,
		SMTPPort: emailCfg.SMTPPort,
		Username: emailCfg.Username,
		Password: emailCfg.Password,
		From:     emailCfg.From,
	})

	// Сервисы
	userSvc := userservice.New(s.UserRepo, s.PasswordResetRepo, emailSender)
	sessionSvc := sessionservice.New(s.SessionRepo)
	chatSvc := chatservice.New(s.ChatRepo, s.SessionRepo, aiClient)
	cardSvc := cardservice.New(s.CardRepo)
	testSvc := testservice.New(s.TestRepo, s.QuestionRepo)
	testCategorySvc := testcategoryservice.New(s.CategoryRepo)
	testViewSvc := testviewservice.New(s.TestViewRepo)
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

	wsHandler := wschat.NewHandler(hub, authService, chatSvc)

	return &App{
		AddChatHandler:              chathandler.Add(log, chatSvc),
		AddMessageByCreatingSession: chathandler.AddByCreatingSession(log, chatSvc),
		GetChatBySessionIdHandler:   chathandler.GetBySessionID(log, chatSvc),
		GetAllSessions:              sessionhandler.GetAll(log, sessionSvc),
		CreateSession:               sessionhandler.Create(log, sessionSvc),
		RetryLastMessage:            chathandler.Retry(log, chatSvc),
		TestCreate:                  testhandler.Create(log, testSvc),
		TestUpdate:                  testhandler.Update(log, testSvc),
		TestGetAll:                  testhandler.GetAll(log, testSvc),
		TestGetById:                 testhandler.GetByID(log, testSvc),
		TestResultsGetAll:           testhandler.GetAllUserResults(log, testSvc),
		TestResultsAdd:              testhandler.AddResult(log, testSvc),
		TestViewAdd:                 testhandler.AddView(log, testViewSvc),
		CreateCategory:              testcategoryhandler.Create(log, testCategorySvc),
		GetAllCategories:            testcategoryhandler.GetAll(log, testCategorySvc),
		CardCreate:                  cardhandler.Create(log, cardSvc),
		CardGetAll:                  cardhandler.GetAll(log, cardSvc),
		CardGetById:                 cardhandler.GetByID(log, cardSvc),
		CardUpdate:                  cardhandler.Update(log, cardSvc),
		UserCreate:                  userhandler.Create(log, userSvc),
		UserLogin:                   userhandler.Login(log, userSvc, jwtService),
		UserForgotPassword:          userhandler.ForgotPassword(log, userSvc),
		UserVerifyCode:              userhandler.VerifyCode(log, userSvc),
		UserResetPassword:           userhandler.ResetPassword(log, userSvc),
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
