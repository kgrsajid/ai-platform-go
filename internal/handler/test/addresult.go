package testhandler

import (
	"fmt"
	"log/slog"
	"net/http"

	"project-go/internal/dto/request"
	res "project-go/internal/dto/response"
	"project-go/internal/lib/auth"
	"project-go/internal/lib/response"
	"project-go/internal/models"
	"project-go/internal/repository/progression"
	testservice "project-go/internal/service/test"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type testResultResponse struct {
	response.Response
	TestResult *res.TestResultResponse `json:"data"`
}

func AddResult(log *slog.Logger, svc *testservice.Service, progRepo *progression.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.test.AddResult"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var testReq request.TestResultReq
		if err := render.DecodeJSON(r.Body, &testReq); err != nil {
			log.Error("failed to decode request")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("invalid request body"))
			return
		}

		if testReq.Score < 0 || testReq.MaxScore < 0 || testReq.TestId == 0 {
			log.Error("invalid request fields")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("invalid request"))
			return
		}

		if testReq.Score > testReq.MaxScore {
			log.Error("score exceeds max score")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("score cannot exceed max score"))
			return
		}

		userId, ok := auth.GetUserID(r)
		if !ok {
			log.Error("unauthorized: missing user id")
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, response.Error("unauthorized"))
			return
		}
		testReq.UserId = userId

		result, err := svc.AddTestResult(testReq)
		if err != nil {
			log.Error("failed to add result", slog.String("error", err.Error()))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to add result"))
			return
		}

		// --- Gamification: award points and XP for quiz completion ---
		go awardQuizGamification(progRepo, userId, testReq.Score, testReq.MaxScore, log)

		render.JSON(w, r, testResultResponse{
			Response:   response.OK(),
			TestResult: result,
		})
	}
}

// awardQuizGamification runs in a goroutine to award points/XP/record activity
// after a quiz is completed. Errors are logged but don't block the response.
func awardQuizGamification(progRepo *progression.Repository, userID uint, score, maxScore int, log *slog.Logger) {
	// Get user's grade band
	_, user, err := progRepo.GetOrCreateProgress(userID)
	if err != nil {
		log.Error("gamification: failed to get user for quiz reward", slog.String("error", err.Error()))
		return
	}

	grade := user.Grade
	if grade == 0 {
		grade = 5
	}
	band := models.GetGradeBand(grade)

	// Calculate percentage
	percentage := float64(0)
	if maxScore > 0 {
		percentage = float64(score) / float64(maxScore) * 100
	}

	// Points and XP based on score percentage
	points := band.QuizPoints(percentage)
	xp := band.QuizXP(percentage)

	if err := progRepo.AddPoints(userID, points, "quiz", fmt.Sprintf("quiz_score_%d", score),
		fmt.Sprintf("Quiz completed: %d/%d (%.0f%%)", score, maxScore, percentage)); err != nil {
		log.Error("gamification: failed to add quiz points", slog.String("error", err.Error()))
	}

	if _, err := progRepo.AddXP(userID, xp, "quiz"); err != nil {
		log.Error("gamification: failed to add quiz XP", slog.String("error", err.Error()))
	}

	// Record activity for streak tracking
	if _, _, err := progRepo.RecordActivity(userID); err != nil {
		log.Error("gamification: failed to record quiz activity", slog.String("error", err.Error()))
	}

	log.Info("gamification: quiz reward awarded",
		slog.Int("userID", int(userID)),
		slog.Int("points", points),
		slog.Int("xp", xp),
		slog.Float64("percentage", percentage),
	)
}
