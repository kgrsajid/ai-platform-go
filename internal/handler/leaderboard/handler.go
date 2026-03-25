package leaderboardhandler

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

// GetLeaderboard returns ranked users
func (h *Handler) GetLeaderboard(w http.ResponseWriter, r *http.Request) {
	const op = "handler.leaderboard.GetLeaderboard"
	log := h.log.With(slog.String("op", op), slog.String("request_id", middleware.GetReqID(r.Context())))

	userID, ok := auth.GetUserID(r)
	if !ok {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, response.Error("unauthorized"))
		return
	}

	// Get query params
	sort := r.URL.Query().Get("sort")
	if sort == "" {
		sort = "level"
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 50 {
		limit = 20
	}

	// Get grade filter (optional)
	var gradeFilter *int
	gradeStr := r.URL.Query().Get("grade")
	if gradeStr != "" {
		if g, err := strconv.Atoi(gradeStr); err == nil && g >= 0 && g <= 11 {
			gradeFilter = &g
		}
	}

	results, total, err := h.repo.GetLeaderboard(sort, page, limit, gradeFilter)
	if err != nil {
		log.Error("failed to get leaderboard", slog.String("error", err.Error()))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, response.Error("failed to get leaderboard"))
		return
	}

	// Mark current user
	for i, entry := range results {
		if entry["user_id"] == userID {
			results[i]["is_current_user"] = true
		}
	}

	render.JSON(w, r, response.OKWithData(map[string]interface{}{
		"leaderboard": results,
		"total":       total,
		"page":        page,
		"limit":       limit,
		"sort":        sort,
	}))
}
