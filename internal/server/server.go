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

	router.Use(chimiddleware.RequestID)
	router.Use(chimiddleware.Logger)
	router.Use(middleware.Logger(log))
	router.Use(chimiddleware.Recoverer)
	router.Use(chimiddleware.URLFormat)

	router.Group(func(r chi.Router) {
		r.Use(authMiddleware)
		r.Post("/chat", app.AddChatHandler)
		r.Put("/chat/retry/{sessionId}", app.RetryLastMessage)
		r.Post("/chat/new", app.AddMessageByCreatingSession)
		r.Post("/session", app.CreateSession)
		r.Get("/session", app.GetAllSessions)
		r.Delete("/session/{sessionId}", app.DeleteSession)
		r.Get("/session/{sessionId}", app.GetChatBySessionIdHandler)
		r.Post("/test", app.TestCreate)
		r.Put("/test/{testId}", app.TestUpdate)
		r.Get("/test", app.TestGetAll)
		r.Post("/test/result", app.TestResultsAdd)
		r.Get("/test/result/{testId}", app.TestResultsGetAll)
		r.Get("/test/{testId}", app.TestGetById)
		r.Post("/test/category", app.CreateCategory)
		r.Get("/test/category", app.GetAllCategories)
		r.Post("/test/view", app.TestViewAdd)
		r.Get("/card", app.CardGetAll)
		r.Post("/card", app.CardCreate)
		r.Get("/card/{cardId}", app.CardGetById)
		r.Put("/card/{cardId}", app.CardUpdate)
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
