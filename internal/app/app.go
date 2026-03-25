package app

import (
	"log/slog"
	"net/http"

	client "project-go/internal/client/chat"
	"project-go/internal/config"
	cardhandler "project-go/internal/handler/card"
	chathandler "project-go/internal/handler/chat"
	leaderboardhandler "project-go/internal/handler/leaderboard"
	progressionhandler "project-go/internal/handler/progression"
	sessionhandler "project-go/internal/handler/session"
	testhandler "project-go/internal/handler/test"
	testcategoryhandler "project-go/internal/handler/testcategory"
	trainerhandler "project-go/internal/handler/trainer"
	userhandler "project-go/internal/handler/user"
	"project-go/internal/lib/auth"
	emaillib "project-go/internal/lib/email"
	"project-go/internal/lib/jwt"
	"project-go/internal/repository/store"
	trainerrepo "project-go/internal/repository/trainer"
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
	GenerateTestDetails         http.HandlerFunc
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
	GenerateCards               http.HandlerFunc
	GetProgression              http.HandlerFunc
	GetStreak                   http.HandlerFunc
	ClaimDailyBonus             http.HandlerFunc
	GetTransactions             http.HandlerFunc
	GetRewards                  http.HandlerFunc
	RedeemReward                http.HandlerFunc
	GetMyRedemptions            http.HandlerFunc
	GetSubjects                 http.HandlerFunc
	GetTrainerProfile           http.HandlerFunc
	UpdateTrainerProfile        http.HandlerFunc
	GetTrainerTimeline          http.HandlerFunc
	GetLeaderboard              http.HandlerFunc

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

	userSvc := userservice.New(s.UserRepo, s.PasswordResetRepo, emailSender)
	sessionSvc := sessionservice.New(s.SessionRepo, aiClient)
	chatSvc := chatservice.New(s.ChatRepo, s.SessionRepo, aiClient)
	cardSvc := cardservice.New(s.CardRepo, aiClient)
	testSvc := testservice.New(s.TestRepo, s.QuestionRepo, aiClient)
	testCategorySvc := testcategoryservice.New(s.CategoryRepo)
	testViewSvc := testviewservice.New(s.TestViewRepo)

	wsHandler := wschat.NewHandler(hub, authService, chatSvc)

	progressionHandler := progressionhandler.New(log, s.ProgressionRepo)
	trainerRepo := trainerrepo.NewRepository(s.DB)
	trainerHandler := trainerhandler.New(log, trainerRepo)
	leaderboardHandler := leaderboardhandler.New(log, trainerRepo)

	return &App{
		AddChatHandler:              chathandler.Add(log, chatSvc),
		AddMessageByCreatingSession: chathandler.AddByCreatingSession(log, chatSvc),
		GetChatBySessionIdHandler:   chathandler.GetBySessionID(log, chatSvc),
		GetAllSessions:              sessionhandler.GetAll(log, sessionSvc),
		CreateSession:               sessionhandler.Create(log, sessionSvc),
		DeleteSession:               sessionhandler.Delete(log, sessionSvc),
		RetryLastMessage:            chathandler.Retry(log, chatSvc),
		TestCreate:                  testhandler.Create(log, testSvc),
		TestUpdate:                  testhandler.Update(log, testSvc),
		TestGetAll:                  testhandler.GetAll(log, testSvc),
		TestGetById:                 testhandler.GetByID(log, testSvc),
		TestResultsGetAll:           testhandler.GetAllUserResults(log, testSvc),
		TestResultsAdd:              testhandler.AddResult(log, testSvc, s.ProgressionRepo),
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
		WSAddMessage:                wsHandler,
		GenerateTestDetails:         testhandler.GenerateQuiz(log, testSvc),
		GenerateCards:               cardhandler.GenerateCard(log, cardSvc),
		GetProgression:              progressionHandler.GetProgression,
		GetStreak:                   progressionHandler.GetStreak,
		ClaimDailyBonus:             progressionHandler.ClaimDailyBonus,
		GetTransactions:             progressionHandler.GetTransactions,
		GetRewards:                  progressionHandler.GetRewards,
		RedeemReward:                progressionHandler.RedeemReward,
		GetMyRedemptions:            progressionHandler.GetMyRedemptions,
		GetSubjects:                 progressionHandler.GetSubjects,
		GetTrainerProfile:           trainerHandler.GetTrainerProfile,
		UpdateTrainerProfile:        trainerHandler.UpdateTrainerProfile,
		GetTrainerTimeline:          trainerHandler.GetTrainerTimeline,
		GetLeaderboard:              leaderboardHandler.GetLeaderboard,
	}
}
