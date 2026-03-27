package assignmentservice

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	client "project-go/internal/client/chat"
	"project-go/internal/models"

	"gorm.io/gorm"
)

type AssignmentRepository interface {
	CreateAssignment(a *models.Assignment) (*models.Assignment, error)
	GetAssignmentByID(id uint) (*models.Assignment, error)
	GetPublishedAssignments(grade int, page, limit int) ([]models.Assignment, int64, error)
	GetTeacherAssignments(teacherID uint, page, limit int) ([]models.Assignment, int64, error)
	UpdateAssignment(a *models.Assignment) (*models.Assignment, error)
	DeleteAssignment(id uint) error
	CreateSubmission(s *models.AssignmentSubmission) (*models.AssignmentSubmission, error)
	GetSubmissionByID(id uint) (*models.AssignmentSubmission, error)
	GetStudentSubmissions(studentID uint, page, limit int) ([]models.AssignmentSubmission, int64, error)
	GetAssignmentSubmissions(assignmentID uint, page, limit int) ([]models.AssignmentSubmission, int64, error)
	UpdateSubmissionAfterEval(s *models.AssignmentSubmission) (*models.AssignmentSubmission, error)
	GetExistingSubmission(studentID, assignmentID uint) (*models.AssignmentSubmission, error)
	GetStudentAssignmentStats(studentID uint) (int64, float64, error)
}

type ProgressionRepository interface {
	AddXP(userID uint, xpAmount int, source string) (*models.UserProgress, error)
	AddPoints(userID uint, amount int, source, referenceID, description string) error
	RecordActivity(userID uint) (*models.DailyActivity, *models.UserProgress, error)
	GetOrCreateProgress(userID uint) (*models.UserProgress, *models.User, error)
}

type Service struct {
	assignRepo AssignmentRepository
	progRepo   ProgressionRepository
	aiAPI      client.AIClient
	log        *slog.Logger
}

func New(assignRepo AssignmentRepository, progRepo ProgressionRepository, aiAPI client.AIClient, log *slog.Logger) *Service {
	return &Service{
		assignRepo: assignRepo,
		progRepo:   progRepo,
		aiAPI:      aiAPI,
		log:        log,
	}
}

// --- Assignment CRUD ---

type CreateAssignmentReq struct {
	Title       string `json:"title" validate:"required"`
	Question    string `json:"question" validate:"required"`
	Rubric      string `json:"rubric"`
	Subject     string `json:"subject"`
	GradeMin    int    `json:"grade_min"`
	GradeMax    int    `json:"grade_max"`
	IsPublished bool   `json:"is_published"`
}

func (s *Service) CreateAssignment(ctx context.Context, teacherID uint, req CreateAssignmentReq) (*models.Assignment, error) {
	a := &models.Assignment{
		TeacherID:   teacherID,
		Title:       req.Title,
		Question:    req.Question,
		Rubric:      req.Rubric,
		Subject:     req.Subject,
		GradeMin:    req.GradeMin,
		GradeMax:    req.GradeMax,
		IsPublished: req.IsPublished,
	}
	return s.assignRepo.CreateAssignment(a)
}

func (s *Service) GetPublishedAssignments(grade int, page, limit int) ([]models.Assignment, int64, error) {
	return s.assignRepo.GetPublishedAssignments(grade, page, limit)
}

func (s *Service) GetTeacherAssignments(teacherID uint, page, limit int) ([]models.Assignment, int64, error) {
	return s.assignRepo.GetTeacherAssignments(teacherID, page, limit)
}

func (s *Service) GetAssignmentByID(id uint) (*models.Assignment, error) {
	return s.assignRepo.GetAssignmentByID(id)
}

func (s *Service) UpdateAssignment(ctx context.Context, teacherID, id uint, req CreateAssignmentReq) (*models.Assignment, error) {
	a, err := s.assignRepo.GetAssignmentByID(id)
	if err != nil {
		return nil, fmt.Errorf("assignment not found: %w", err)
	}
	if a.TeacherID != teacherID {
		return nil, fmt.Errorf("not authorized to edit this assignment")
	}

	a.Title = req.Title
	a.Question = req.Question
	a.Rubric = req.Rubric
	a.Subject = req.Subject
	a.GradeMin = req.GradeMin
	a.GradeMax = req.GradeMax
	a.IsPublished = req.IsPublished

	return s.assignRepo.UpdateAssignment(a)
}

func (s *Service) DeleteAssignment(teacherID, id uint) error {
	a, err := s.assignRepo.GetAssignmentByID(id)
	if err != nil {
		return fmt.Errorf("assignment not found: %w", err)
	}
	if a.TeacherID != teacherID {
		return fmt.Errorf("not authorized to delete this assignment")
	}
	return s.assignRepo.DeleteAssignment(id)
}

// --- Submissions ---

func (s *Service) SubmitAnswer(ctx context.Context, studentID, assignmentID uint, answer string) (*models.AssignmentSubmission, error) {
	// Check if already submitted
	existing, err := s.assignRepo.GetExistingSubmission(studentID, assignmentID)
	if err == nil && existing != nil {
		return nil, fmt.Errorf("you have already submitted an answer for this assignment")
	}

	// Get assignment
	assignment, err := s.assignRepo.GetAssignmentByID(assignmentID)
	if err != nil {
		return nil, fmt.Errorf("assignment not found: %w", err)
	}

	// Create submission
	submission := &models.AssignmentSubmission{
		StudentID:    studentID,
		AssignmentID: assignmentID,
		Answer:       answer,
		IsEvaluated:  false,
	}

	submission, err = s.assignRepo.CreateSubmission(submission)
	if err != nil {
		return nil, err
	}

	// Evaluate with AI (synchronous — student waits for result)
	if s.aiAPI != nil {
		grade := s.getUserGrade(studentID)
		evalResult, err := s.aiAPI.EvaluateAssignment(ctx, assignment.Question, answer, assignment.Rubric, grade)
		if err != nil {
			s.log.Error("assignment evaluation failed", slog.String("error", err.Error()))
			// Return submission without evaluation — can be retried
			return submission, nil
		}

		// Update submission with evaluation
		now := time.Now()
		submission.Score = evalResult.Score
		submission.MaxScore = evalResult.MaxScore
		submission.Feedback = evalResult.Feedback
		submission.GradeLevel = evalResult.GradeLevel
		submission.IsEvaluated = true
		submission.EvaluatedAt = &now

		strengthsJSON, _ := json.Marshal(evalResult.Strengths)
		improvementsJSON, _ := json.Marshal(evalResult.Improvements)
		submission.Strengths = string(strengthsJSON)
		submission.Improvements = string(improvementsJSON)

		submission, err = s.assignRepo.UpdateSubmissionAfterEval(submission)
		if err != nil {
			s.log.Error("failed to save evaluation", slog.String("error", err.Error()))
		}

		// Award gamification based on score
		go s.awardAssignmentGamification(studentID, evalResult.Score, uint(assignmentID))
	}

	return submission, nil
}

// awardAssignmentGamification awards XP and points based on assignment score
func (s *Service) awardAssignmentGamification(userID uint, score int, assignmentID uint) {
	// XP based on score
	xp := score / 10 // 0-10 XP based on score percentage
	if xp < 1 && score > 0 {
		xp = 1
	}

	// Bonus XP for high scores
	if score >= 90 {
		xp += 5
	} else if score >= 70 {
		xp += 3
	} else if score >= 50 {
		xp += 1
	}

	if _, err := s.progRepo.AddXP(userID, xp, "assignment"); err != nil {
		s.log.Error("gamification: failed to add assignment XP", slog.String("error", err.Error()))
	}

	// Points based on score (1 point per 10 score, minimum 1)
	points := score / 10
	if points < 1 && score > 0 {
		points = 1
	}
	if err := s.progRepo.AddPoints(userID, points, "assignment", fmt.Sprintf("%d", assignmentID), fmt.Sprintf("Assignment score: %d/100", score)); err != nil {
		s.log.Error("gamification: failed to add assignment points", slog.String("error", err.Error()))
	}

	// Record activity
	if _, _, err := s.progRepo.RecordActivity(userID); err != nil {
		s.log.Error("gamification: failed to record assignment activity", slog.String("error", err.Error()))
	}

	s.log.Info("gamification: assignment reward awarded",
		slog.Int("userID", int(userID)),
		slog.Int("score", score),
		slog.Int("xp", xp),
		slog.Int("points", points),
	)
}

func (s *Service) getUserGrade(userID uint) int {
	_, user, err := s.progRepo.GetOrCreateProgress(userID)
	if err != nil || user == nil {
		return 5
	}
	if user.Grade <= 0 || user.Grade > 11 {
		return 5
	}
	return user.Grade
}

func (s *Service) GetStudentSubmissions(studentID uint, page, limit int) ([]models.AssignmentSubmission, int64, error) {
	return s.assignRepo.GetStudentSubmissions(studentID, page, limit)
}

func (s *Service) GetAssignmentSubmissions(teacherID, assignmentID uint, page, limit int) ([]models.AssignmentSubmission, int64, error) {
	// Verify teacher owns this assignment
	assignment, err := s.assignRepo.GetAssignmentByID(assignmentID)
	if err != nil {
		return nil, 0, fmt.Errorf("assignment not found: %w", err)
	}
	if assignment.TeacherID != teacherID {
		return nil, 0, fmt.Errorf("not authorized")
	}
	return s.assignRepo.GetAssignmentSubmissions(assignmentID, page, limit)
}

func (s *Service) GetStudentStats(studentID uint) (totalSubmitted int64, avgScore float64, err error) {
	return s.assignRepo.GetStudentAssignmentStats(studentID)
}
