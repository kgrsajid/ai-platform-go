package server

import (
	"log/slog"
	"net/http"

	"project-go/internal/app"
	"project-go/internal/middleware"
	"project-go/internal/repository/store"

	chimiddleware "github.com/go-chi/chi/v5/middleware"

	"github.com/go-chi/chi/v5"
)

func NewRouter(app *app.App, log *slog.Logger, s *store.Store, jwtKey string) http.Handler {
	authMiddleware := middleware.AuthMiddleware([]byte(jwtKey))

	router := chi.NewRouter()

	router.Use(middleware.CORSMiddleware)
	router.Use(chimiddleware.RequestID)
	router.Use(chimiddleware.Logger)
	router.Use(middleware.Logger(log))
	router.Use(chimiddleware.Recoverer)

	router.Group(func(r chi.Router) {
		r.Use(authMiddleware)
		r.Post("/chat", app.AddChatHandler)
		r.Put("/chat/retry/{sessionId}", app.RetryLastMessage)
		r.Post("/chat/new", app.AddMessageByCreatingSession)
		r.Post("/session", app.CreateSession)
		r.Get("/session", app.GetAllSessions)
		r.Delete("/session/{sessionId}", app.DeleteSession)
		r.Get("/session/{sessionId}", app.GetChatBySessionIdHandler)
		r.Route("/test", func(r chi.Router) {
			r.Post("/", app.TestCreate)
			r.Put("/{testId}", app.TestUpdate)
			r.Get("/", app.TestGetAll)
			r.Post("/result", app.TestResultsAdd)
			r.Get("/result/{testId}", app.TestResultsGetAll)
			r.Get("/{testId}", app.TestGetById)
			r.Post("/generate", app.GenerateTestDetails)
			r.Post("/category", app.CreateCategory)
			r.Get("/category", app.GetAllCategories)
			r.Post("/view", app.TestViewAdd)
		})
		r.Route("/card", func(r chi.Router) {
			r.Get("/", app.CardGetAll)
			r.Post("/", app.CardCreate)
			r.Get("/{cardId}", app.CardGetById)
			r.Put("/{cardId}", app.CardUpdate)
			r.Post("/generate", app.GenerateCards)
		})
		// Gamification endpoints (Phase 0)
		r.Get("/progression", app.GetProgression)
		r.Get("/progression/streak", app.GetStreak)
		r.Post("/progression/streak/claim", app.ClaimDailyBonus)
		r.Get("/progression/transactions", app.GetTransactions)
		r.Get("/rewards", app.GetRewards)
		r.Post("/rewards/{id}/redeem", app.RedeemReward)
		r.Get("/rewards/my", app.GetMyRedemptions)
		r.Get("/subjects", app.GetSubjects)
	})

	router.Get("/message", app.WSAddMessage.ServeWS)
	router.Route("/auth", func(r chi.Router) {
		r.Post("/register", app.UserCreate)
		r.Post("/login", app.UserLogin)
		r.Post("/forgot-password", app.UserForgotPassword)
		r.Post("/verify-code", app.UserVerifyCode)
		r.Post("/reset-password", app.UserResetPassword)
	})

	return router
}
