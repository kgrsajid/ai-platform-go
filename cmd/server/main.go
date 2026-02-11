package main

import (
	"log/slog"
	"net/http"
	"os"
	"project-go/internal/app"
	"project-go/internal/config"
	"project-go/internal/http-server/middleware/cors"
	"project-go/internal/http-server/repository/store"
	"project-go/internal/lib/logger/sl"
	"project-go/internal/logger"
	"project-go/internal/server"
	"project-go/internal/storage/postgres"
)

func main() {
	cfg := config.MustLoad()

	log := logger.SetupLogger(cfg.Env)
	log = log.With("env", cfg.Env)
	db, err := postgres.New(cfg.Dsn)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}
	store := store.NewStore(db)
	app := app.New(log, store, cfg.JWT_Key, cfg.AI_Base_Url)
	router := server.NewRouter(app, log, store, cfg.JWT_Key)
	handler := cors.CORSMiddleware(router)
	log.Info("starting server", slog.String("address", cfg.Address))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      handler,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start and run server")
	}

	log.Error("server stopped")

}
