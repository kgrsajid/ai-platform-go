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
		r.Post("/chat", app.CreateChatHandler)
		r.Get("/session", app.GetAllSessions)
		r.Get("/session/{sessionId}", app.GetChatBySessionIdHandler)
		r.Post("/test", app.TestCreate)
	})
	router.Route("/auth", func(r chi.Router) {
		r.Post("/register", app.UserCreate)
		r.Post("/login", app.UserLogin)
	})

	return router
}
