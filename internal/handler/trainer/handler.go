package trainerhandler

import (
	"log/slog"
	"net/http"
	"strconv"

	"project-go/internal/lib/auth"
	"project-go/internal/lib/response"
	"project-go/internal/repository/trainer"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Handler struct {
	log  *slog.Logger
	repo *trainer.Repository
}

func New(log *slog.Logger, repo *trainer.Repository) *Handler {
	return &Handler{log: log, repo: repo}
}

// GetTrainerProfile returns the robot profile with computed stats
func (h *Handler) GetTrainerProfile(w http.ResponseWriter, r *http.Request) {
	const op = "handler.trainer.GetTrainerProfile"
	log := h.log.With(slog.String("op", op), slog.String("request_id", middleware.GetReqID(r.Context())))

	userID, ok := auth.GetUserID(r)
	if !ok {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, response.Error("unauthorized"))
		return
	}

	progress, user, err := h.repo.GetTrainerProfile(userID)
	if err != nil {
		log.Error("failed to get trainer profile", slog.String("error", err.Error()))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, response.Error("failed to get trainer profile"))
		return
	}

	stats, err := h.repo.GetTrainerStats(userID)
	if err != nil {
		log.Error("failed to get trainer stats", slog.String("error", err.Error()))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, response.Error("failed to get trainer stats"))
		return
	}

	data := map[string]interface{}{
		"robot_name":    progress.RobotName,
		"robot_color":   progress.RobotColor,
		"robot_level":   progress.RobotLevel,
		"current_level": progress.CurrentLevel,
		"stage_name":    stats["stage_name"],
		"total_xp":      progress.TotalXP,
		"grade":         user.Grade,
		"stats":         stats,
		"created_at":    progress.CreatedAt,
	}

	render.JSON(w, r, response.OKWithData(data))
}

// UpdateTrainerProfile updates robot name and color
func (h *Handler) UpdateTrainerProfile(w http.ResponseWriter, r *http.Request) {
	const op = "handler.trainer.UpdateTrainerProfile"
	log := h.log.With(slog.String("op", op), slog.String("request_id", middleware.GetReqID(r.Context())))

	userID, ok := auth.GetUserID(r)
	if !ok {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, response.Error("unauthorized"))
		return
	}

	var req struct {
		Name  string `json:"name"`
		Color string `json:"color"`
	}

	if err := render.DecodeJSON(r.Body, &req); err != nil {
		log.Error("failed to decode request", slog.String("error", err.Error()))
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Error("invalid request body"))
		return
	}

	progress, err := h.repo.UpdateTrainerProfile(userID, req.Name, req.Color)
	if err != nil {
		log.Error("failed to update trainer profile", slog.String("error", err.Error()))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, response.Error("failed to update trainer profile"))
		return
	}

	data := map[string]interface{}{
		"robot_name":  progress.RobotName,
		"robot_color": progress.RobotColor,
	}

	render.JSON(w, r, response.OKWithData(data))
}

// GetTrainerTimeline returns level-up history
func (h *Handler) GetTrainerTimeline(w http.ResponseWriter, r *http.Request) {
	const op = "handler.trainer.GetTrainerTimeline"
	log := h.log.With(slog.String("op", op), slog.String("request_id", middleware.GetReqID(r.Context())))

	userID, ok := auth.GetUserID(r)
	if !ok {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, response.Error("unauthorized"))
		return
	}

	limitStr := r.URL.Query().Get("limit")
	limit := 20
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 50 {
			limit = l
		}
	}

	transactions, err := h.repo.GetTrainerTimeline(userID, limit)
	if err != nil {
		log.Error("failed to get trainer timeline", slog.String("error", err.Error()))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, response.Error("failed to get trainer timeline"))
		return
	}

	render.JSON(w, r, response.OKWithData(map[string]interface{}{
		"timeline": transactions,
	}))
}
