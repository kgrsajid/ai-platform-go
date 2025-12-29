package server

import (
	"log/slog"
	"net/http"
	"project-go/internal/app"
	"project-go/internal/http-server/middleware/auth"
	"project-go/internal/http-server/middleware/logger"
	"project-go/internal/http-server/repository/store"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(app *app.App, log *slog.Logger, store *store.Store, jwtKey string) http.Handler {
	authMiddleware := auth.AuthMiddleware([]byte(jwtKey))

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(logger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Group(func(r chi.Router) {
		r.Use(authMiddleware)
		r.Post("/chat", app.AddChatHandler)
		r.Post("/chat/new", app.AddMessageByCreatingSession)
		r.Get("/session", app.GetAllSessions)
		r.Get("/session/{sessionId}", app.GetChatBySessionIdHandler)
		r.Post("/test", app.TestCreate)
		r.Get("/test", app.TestGetAll)
		r.Post("/test/result", app.TestResultsAdd)
		r.Get("/test/result/{testId}", app.TestResultsGetALl)
		r.Get("/test/{testId}", app.TestGetById)
		r.Post("/test/category", app.CreateCategory)
		r.Get("/test/category", app.GetAllCategories)
		r.Post("/test/view", app.TestViewAdd)
	})
	router.Route("/auth", func(r chi.Router) {
		r.Post("/register", app.UserCreate)
		r.Post("/login", app.UserLogin)
	})

	return router
}
