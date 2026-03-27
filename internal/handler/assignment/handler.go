package assignmenthandler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"project-go/internal/lib/auth"
	"project-go/internal/lib/response"
	"project-go/internal/models"
	assignmentservice "project-go/internal/service/assignment"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type listResponse struct {
	response.Response
	Data  interface{} `json:"data"`
	Total int64       `json:"total"`
}

type detailResponse struct {
	response.Response
	Data interface{} `json:"data"`
}

func New(log *slog.Logger, svc *assignmentservice.Service) *Handler {
	return &Handler{log: log, svc: svc}
}

type Handler struct {
	log *slog.Logger
	svc *assignmentservice.Service
}

// CreateAssignment — POST /assignment (teacher only)
func (h *Handler) CreateAssignment(w http.ResponseWriter, r *http.Request) {
	const op = "handler.assignment.Create"
	log := h.log.With(slog.String("op", op), slog.String("request_id", middleware.GetReqID(r.Context())))

	var req assignmentservice.CreateAssignmentReq
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		log.Error("failed to decode request")
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Error("invalid request body"))
		return
	}

	if req.Title == "" || req.Question == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Error("title and question are required"))
		return
	}

	userID, ok := auth.GetUserID(r)
	if !ok {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, response.Error("unauthorized"))
		return
	}

	assignment, err := h.svc.CreateAssignment(r.Context(), userID, req)
	if err != nil {
		log.Error("failed to create assignment", slog.String("error", err.Error()))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, response.Error("failed to create assignment: "+err.Error()))
		return
	}

	render.JSON(w, r, detailResponse{Response: response.OK(), Data: toAssignmentItem(assignment)})
}

// GetPublishedAssignments — GET /assignment (student view)
func (h *Handler) GetPublishedAssignments(w http.ResponseWriter, r *http.Request) {
	log := h.log.With(slog.String("op", "handler.assignment.GetPublished"))

	grade, _ := strconv.Atoi(r.URL.Query().Get("grade"))
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if grade <= 0 { grade = 5 }
	if page <= 0 { page = 1 }
	if limit <= 0 || limit > 50 { limit = 20 }

	assignments, total, err := h.svc.GetPublishedAssignments(grade, page, limit)
	if err != nil {
		log.Error("failed to get assignments", slog.String("error", err.Error()))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, response.Error("failed to get assignments"))
		return
	}

	items := make([]map[string]interface{}, len(assignments))
	for i, a := range assignments {
		items[i] = toAssignmentItem(&a)
	}

	render.JSON(w, r, listResponse{Response: response.OK(), Data: items, Total: total})
}

// GetMyAssignments — GET /assignment/my (teacher view)
func (h *Handler) GetMyAssignments(w http.ResponseWriter, r *http.Request) {
	log := h.log.With(slog.String("op", "handler.assignment.GetMy"))

	userID, ok := auth.GetUserID(r)
	if !ok {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, response.Error("unauthorized"))
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if page <= 0 { page = 1 }
	if limit <= 0 || limit > 50 { limit = 20 }

	assignments, total, err := h.svc.GetTeacherAssignments(userID, page, limit)
	if err != nil {
		log.Error("failed to get assignments", slog.String("error", err.Error()))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, response.Error("failed to get assignments"))
		return
	}

	items := make([]map[string]interface{}, len(assignments))
	for i, a := range assignments {
		items[i] = toAssignmentItem(&a)
	}

	render.JSON(w, r, listResponse{Response: response.OK(), Data: items, Total: total})
}

// GetAssignmentByID — GET /assignment/{id}
func (h *Handler) GetAssignmentByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Error("invalid assignment id"))
		return
	}

	assignment, err := h.svc.GetAssignmentByID(uint(id))
	if err != nil {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, response.Error("assignment not found"))
		return
	}

	render.JSON(w, r, detailResponse{Response: response.OK(), Data: toAssignmentItem(assignment)})
}

// UpdateAssignment — PUT /assignment/{id}
func (h *Handler) UpdateAssignment(w http.ResponseWriter, r *http.Request) {
	log := h.log.With(slog.String("op", "handler.assignment.Update"))

	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Error("invalid assignment id"))
		return
	}

	var req assignmentservice.CreateAssignmentReq
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Error("invalid request body"))
		return
	}

	userID, ok := auth.GetUserID(r)
	if !ok {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, response.Error("unauthorized"))
		return
	}

	assignment, err := h.svc.UpdateAssignment(r.Context(), userID, uint(id), req)
	if err != nil {
		log.Error("failed to update assignment", slog.String("error", err.Error()))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, response.Error(err.Error()))
		return
	}

	render.JSON(w, r, detailResponse{Response: response.OK(), Data: toAssignmentItem(assignment)})
}

// DeleteAssignment — DELETE /assignment/{id}
func (h *Handler) DeleteAssignment(w http.ResponseWriter, r *http.Request) {
	log := h.log.With(slog.String("op", "handler.assignment.Delete"))

	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Error("invalid assignment id"))
		return
	}

	userID, ok := auth.GetUserID(r)
	if !ok {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, response.Error("unauthorized"))
		return
	}

	if err := h.svc.DeleteAssignment(userID, uint(id)); err != nil {
		log.Error("failed to delete assignment", slog.String("error", err.Error()))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, response.Error(err.Error()))
		return
	}

	render.JSON(w, r, response.OK())
}

// SubmitAnswer — POST /assignment/{id}/submit
func (h *Handler) SubmitAnswer(w http.ResponseWriter, r *http.Request) {
	log := h.log.With(slog.String("op", "handler.assignment.Submit"))

	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Error("invalid assignment id"))
		return
	}

	var body struct {
		Answer string `json:"answer"`
	}
	if err := render.DecodeJSON(r.Body, &body); err != nil || body.Answer == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Error("answer is required"))
		return
	}

	userID, ok := auth.GetUserID(r)
	if !ok {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, response.Error("unauthorized"))
		return
	}

	submission, err := h.svc.SubmitAnswer(r.Context(), userID, uint(id), body.Answer)
	if err != nil {
		log.Error("failed to submit answer", slog.String("error", err.Error()))
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Error(err.Error()))
		return
	}

	render.JSON(w, r, detailResponse{Response: response.OK(), Data: toSubmissionItem(submission)})
}

// GetMySubmissions — GET /assignment/submissions/my (student view)
func (h *Handler) GetMySubmissions(w http.ResponseWriter, r *http.Request) {
	log := h.log.With(slog.String("op", "handler.assignment.GetMySubmissions"))

	userID, ok := auth.GetUserID(r)
	if !ok {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, response.Error("unauthorized"))
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if page <= 0 { page = 1 }
	if limit <= 0 || limit > 50 { limit = 20 }

	submissions, total, err := h.svc.GetStudentSubmissions(userID, page, limit)
	if err != nil {
		log.Error("failed to get submissions", slog.String("error", err.Error()))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, response.Error("failed to get submissions"))
		return
	}

	items := make([]map[string]interface{}, len(submissions))
	for i, s := range submissions {
		items[i] = toSubmissionItem(&s)
	}

	render.JSON(w, r, listResponse{Response: response.OK(), Data: items, Total: total})
}

// GetAssignmentSubmissions — GET /assignment/{id}/submissions (teacher view)
func (h *Handler) GetAssignmentSubmissions(w http.ResponseWriter, r *http.Request) {
	log := h.log.With(slog.String("op", "handler.assignment.GetSubmissions"))

	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Error("invalid assignment id"))
		return
	}

	userID, ok := auth.GetUserID(r)
	if !ok {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, response.Error("unauthorized"))
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if page <= 0 { page = 1 }
	if limit <= 0 || limit > 50 { limit = 20 }

	submissions, total, err := h.svc.GetAssignmentSubmissions(userID, uint(id), page, limit)
	if err != nil {
		log.Error("failed to get submissions", slog.String("error", err.Error()))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, response.Error(err.Error()))
		return
	}

	items := make([]map[string]interface{}, len(submissions))
	for i, s := range submissions {
		items[i] = toSubmissionItem(&s)
	}

	render.JSON(w, r, listResponse{Response: response.OK(), Data: items, Total: total})
}

// GetStudentStats — GET /assignment/stats/my
func (h *Handler) GetStudentStats(w http.ResponseWriter, r *http.Request) {
	log := h.log.With(slog.String("op", "handler.assignment.GetStudentStats"))

	userID, ok := auth.GetUserID(r)
	if !ok {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, response.Error("unauthorized"))
		return
	}

	total, avg, err := h.svc.GetStudentStats(userID)
	if err != nil {
		log.Error("failed to get stats", slog.String("error", err.Error()))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, response.Error("failed to get stats"))
		return
	}

	render.JSON(w, r, detailResponse{
		Response: response.OK(),
		Data: map[string]interface{}{
			"total_submitted": total,
			"average_score":   avg,
		},
	})
}

// --- Helpers ---

func toAssignmentItem(a *models.Assignment) map[string]interface{} {
	teacherName := ""
	if a.Teacher.Name != "" {
		teacherName = a.Teacher.Name
	}
	return map[string]interface{}{
		"id":           a.ID,
		"teacher_id":   a.TeacherID,
		"title":        a.Title,
		"question":     a.Question,
		"rubric":       a.Rubric,
		"subject":      a.Subject,
		"grade_min":    a.GradeMin,
		"grade_max":    a.GradeMax,
		"is_published": a.IsPublished,
		"teacher_name": teacherName,
		"created_at":   a.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

func toSubmissionItem(s *models.AssignmentSubmission) map[string]interface{} {
	var strengths, improvements []string
	json.Unmarshal([]byte(s.Strengths), &strengths)
	json.Unmarshal([]byte(s.Improvements), &improvements)

	studentName := ""
	if s.Student.Name != "" {
		studentName = s.Student.Name
	}
	title := ""
	if s.Assignment.Title != "" {
		title = s.Assignment.Title
	}

	return map[string]interface{}{
		"id":             s.ID,
		"assignment_id":  s.AssignmentID,
		"answer":         s.Answer,
		"score":          s.Score,
		"max_score":      s.MaxScore,
		"feedback":       s.Feedback,
		"strengths":      strengths,
		"improvements":   improvements,
		"grade_level":    s.GradeLevel,
		"is_evaluated":   s.IsEvaluated,
		"student_name":   studentName,
		"title":          title,
		"created_at":     s.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
