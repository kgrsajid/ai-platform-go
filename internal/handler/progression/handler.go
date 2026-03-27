package progressionhandler

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"project-go/internal/lib/auth"
	"project-go/internal/lib/response"
	"project-go/internal/models"
	"project-go/internal/repository/progression"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type Handler struct {
	log *slog.Logger
	repo *progression.Repository
}

func New(log *slog.Logger, repo *progression.Repository) *Handler {
	return &Handler{log: log, repo: repo}
}

type progressionResponse struct {
	response.Response
	Data *progressionData `json:"data"`
}

type progressionData struct {
	TotalPoints     int    `json:"total_points"`
	AvailablePoints int    `json:"available_points"`
	CurrentLevel    int    `json:"current_level"`
	TotalXP         int    `json:"total_xp"`
	CurrentStreak   int    `json:"current_streak"`
	LongestStreak   int    `json:"longest_streak"`
	LevelName       string `json:"level_name"`
	CurrencyIcon    string `json:"currency_icon"`
	CurrencyLabel   string `json:"currency_label"`
	NextLevelXP     int    `json:"next_level_xp"`
	ProgressPercent float64 `json:"progress_percent"`
}

// GetProgression returns the user's full progression profile
func (h *Handler) GetProgression(w http.ResponseWriter, r *http.Request) {
	const op = "handler.progression.GetProgression"
	log := h.log.With(slog.String("op", op), slog.String("request_id", middleware.GetReqID(r.Context())))

	userID, ok := auth.GetUserID(r)
	if !ok {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, response.Error("unauthorized"))
		return
	}

	progress, user, err := h.repo.GetOrCreateProgress(userID)
	if err != nil {
		log.Error("failed to get progression", slog.String("error", err.Error()))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, response.Error("failed to get progression"))
		return
	}

	// Get user's grade for band calculation
	grade := user.Grade
	if grade <= 0 || grade > 11 {
		grade = 5
	}
	band := models.GetGradeBand(grade)
	xpPerLevel := band.XPPerLevel()

	// Calculate progress to next level
	xpForCurrentLevel := 0
	for i := 1; i < progress.CurrentLevel; i++ {
		xpForCurrentLevel += xpPerLevel
	}
	remainingXP := progress.TotalXP - xpForCurrentLevel
	if remainingXP < 0 {
		remainingXP = 0
	}
	progressPercent := float64(remainingXP) / float64(xpPerLevel) * 100
	if progressPercent > 100 {
		progressPercent = 100
	}

	data := &progressionData{
		TotalPoints:     progress.TotalPoints,
		AvailablePoints: progress.AvailablePoints,
		CurrentLevel:    progress.CurrentLevel,
		TotalXP:         progress.TotalXP,
		CurrentStreak:   progress.CurrentStreak,
		LongestStreak:   progress.LongestStreak,
		LevelName:       band.LevelName(progress.CurrentLevel),
		CurrencyIcon:    band.CurrencyIcon(),
		CurrencyLabel:   band.CurrencyLabel(),
		NextLevelXP:     xpPerLevel,
		ProgressPercent: progressPercent,
	}

	render.JSON(w, r, progressionResponse{
		Response: response.OK(),
		Data:     data,
	})
}

// GetStreak returns streak info
func (h *Handler) GetStreak(w http.ResponseWriter, r *http.Request) {
	const op = "handler.progression.GetStreak"
	log := h.log.With(slog.String("op", op), slog.String("request_id", middleware.GetReqID(r.Context())))

	userID, ok := auth.GetUserID(r)
	if !ok {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, response.Error("unauthorized"))
		return
	}

	progress, _, err := h.repo.GetOrCreateProgress(userID)
	if err != nil {
		log.Error("failed to get progression", slog.String("error", err.Error()))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, response.Error("failed to get progression"))
		return
	}

	render.JSON(w, r, response.OKWithData(map[string]interface{}{
		"current_streak": progress.CurrentStreak,
		"longest_streak": progress.LongestStreak,
		"last_active_at": progress.LastActiveAt,
	}))
}

// ClaimDailyBonus records activity and awards daily login bonus
func (h *Handler) ClaimDailyBonus(w http.ResponseWriter, r *http.Request) {
	const op = "handler.progression.ClaimDailyBonus"
	log := h.log.With(slog.String("op", op), slog.String("request_id", middleware.GetReqID(r.Context())))

	userID, ok := auth.GetUserID(r)
	if !ok {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, response.Error("unauthorized"))
		return
	}

	// Get user's grade first
	progress, user, err := h.repo.GetOrCreateProgress(userID)
	if err != nil {
		log.Error("failed to get progression", slog.String("error", err.Error()))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, response.Error("failed to get progression"))
		return
	}

	grade := user.Grade
	if grade <= 0 || grade > 11 {
		grade = 5
	}

	// Record activity (handles streak)
	_, progress, err = h.repo.RecordActivity(userID)
	if err != nil {
		log.Error("failed to record activity", slog.String("error", err.Error()))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, response.Error("failed to record activity"))
		return
	}

	// Award daily login points and XP
	band := models.GetGradeBand(grade)

	pointsEarned := band.DailyLoginPoints()
	xpEarned := band.DailyLoginXP()

	if err := h.repo.AddPoints(userID, pointsEarned, "daily_login", "", "Daily login bonus"); err != nil {
		log.Error("failed to add points", slog.String("error", err.Error()))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, response.Error("failed to add points"))
		return
	}

	if _, err := h.repo.AddXP(userID, xpEarned, "daily_login"); err != nil {
		log.Error("failed to add XP", slog.String("error", err.Error()))
	}

	// Award streak bonus if applicable
	streakBonus := band.StreakBonusPoints(progress.CurrentStreak)
	if streakBonus > 0 {
		if err := h.repo.AddPoints(userID, streakBonus, "streak_bonus", "",
			fmt.Sprintf("%d-day streak bonus!", progress.CurrentStreak)); err != nil {
			log.Error("failed to add streak bonus", slog.String("error", err.Error()))
		}
	}

	render.JSON(w, r, response.OKWithData(map[string]interface{}{
		"points_earned":  pointsEarned,
		"xp_earned":      xpEarned,
		"streak_bonus":   streakBonus,
		"current_streak": progress.CurrentStreak,
		"longest_streak": progress.LongestStreak,
	}))
}

// GetTransactions returns paginated point transactions
func (h *Handler) GetTransactions(w http.ResponseWriter, r *http.Request) {
	const op = "handler.progression.GetTransactions"
	log := h.log.With(slog.String("op", op), slog.String("request_id", middleware.GetReqID(r.Context())))

	userID, ok := auth.GetUserID(r)
	if !ok {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, response.Error("unauthorized"))
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 50 {
		limit = 20
	}

	transactions, total, err := h.repo.GetTransactions(userID, page, limit)
	if err != nil {
		log.Error("failed to get transactions", slog.String("error", err.Error()))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, response.Error("failed to get transactions"))
		return
	}

	render.JSON(w, r, response.OKWithData(map[string]interface{}{
		"transactions": transactions,
		"total":        total,
		"page":         page,
		"limit":        limit,
	}))
}

// GetRewards returns available rewards filtered by user's grade
func (h *Handler) GetRewards(w http.ResponseWriter, r *http.Request) {
	const op = "handler.progression.GetRewards"
	log := h.log.With(slog.String("op", op), slog.String("request_id", middleware.GetReqID(r.Context())))

	userID, ok := auth.GetUserID(r)
	if !ok {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, response.Error("unauthorized"))
		return
	}

	_, user, err := h.repo.GetOrCreateProgress(userID)
	if err != nil {
		log.Error("failed to get progression", slog.String("error", err.Error()))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, response.Error("failed to get progression"))
		return
	}

	grade := user.Grade
	if grade <= 0 || grade > 11 {
		grade = 5
	}
	rewards, err := h.repo.GetRewards(grade)
	if err != nil {
		log.Error("failed to get rewards", slog.String("error", err.Error()))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, response.Error("failed to get rewards"))
		return
	}

	render.JSON(w, r, response.OKWithData(map[string]interface{}{
		"rewards": rewards,
	}))
}

// RedeemReward redeems a reward
func (h *Handler) RedeemReward(w http.ResponseWriter, r *http.Request) {
	const op = "handler.progression.RedeemReward"
	log := h.log.With(slog.String("op", op), slog.String("request_id", middleware.GetReqID(r.Context())))

	userID, ok := auth.GetUserID(r)
	if !ok {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, response.Error("unauthorized"))
		return
	}

	rewardIDStr := chi.URLParam(r, "id")
	rewardID, err := strconv.ParseUint(rewardIDStr, 10, 32)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Error("invalid reward id"))
		return
	}

	redemption, err := h.repo.RedeemReward(userID, uint(rewardID))
	if err != nil {
		log.Error("failed to redeem reward", slog.String("error", err.Error()))
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Error(err.Error()))
		return
	}

	render.JSON(w, r, response.OKWithData(redemption))
}

// GetMyRedemptions returns user's redeemed coupons
func (h *Handler) GetMyRedemptions(w http.ResponseWriter, r *http.Request) {
	const op = "handler.progression.GetMyRedemptions"
	log := h.log.With(slog.String("op", op), slog.String("request_id", middleware.GetReqID(r.Context())))

	userID, ok := auth.GetUserID(r)
	if !ok {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, response.Error("unauthorized"))
		return
	}

	redemptions, err := h.repo.GetMyRedemptions(userID)
	if err != nil {
		log.Error("failed to get redemptions", slog.String("error", err.Error()))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, response.Error("failed to get redemptions"))
		return
	}

	render.JSON(w, r, response.OKWithData(map[string]interface{}{
		"redemptions": redemptions,
	}))
}

// GetSubjects returns all subjects
func (h *Handler) GetSubjects(w http.ResponseWriter, r *http.Request) {
	subjects, err := h.repo.GetSubjects()
	if err != nil {
		h.log.Error("failed to get subjects", slog.String("error", err.Error()))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, response.Error("failed to get subjects"))
		return
	}

	render.JSON(w, r, response.OKWithData(map[string]interface{}{
		"subjects": subjects,
	}))
}
