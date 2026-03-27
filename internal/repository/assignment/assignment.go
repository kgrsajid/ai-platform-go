package assignment

import (
	"project-go/internal/models"
	"time"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// AutoMigrate creates assignment tables
func (r *Repository) AutoMigrate() error {
	return r.db.AutoMigrate(&models.Assignment{}, &models.AssignmentSubmission{})
}

// --- Assignment CRUD ---

func (r *Repository) CreateAssignment(a *models.Assignment) (*models.Assignment, error) {
	if err := r.db.Create(a).Error; err != nil {
		return nil, err
	}
	return a, nil
}

func (r *Repository) GetAssignmentByID(id uint) (*models.Assignment, error) {
	var a models.Assignment
	if err := r.db.Preload("Teacher").First(&a, id).Error; err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *Repository) GetPublishedAssignments(grade int, page, limit int) ([]models.Assignment, int64, error) {
	var assignments []models.Assignment
	var total int64

	query := r.db.Where("is_published = ?", true).
		Where("grade_min <= ? AND grade_max >= ?", grade, grade)

	if err := query.Model(&models.Assignment{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := query.Preload("Teacher").
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&assignments).Error; err != nil {
		return nil, 0, err
	}

	return assignments, total, nil
}

func (r *Repository) GetTeacherAssignments(teacherID uint, page, limit int) ([]models.Assignment, int64, error) {
	var assignments []models.Assignment
	var total int64

	query := r.db.Where("teacher_id = ?", teacherID)

	if err := query.Model(&models.Assignment{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := query.Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&assignments).Error; err != nil {
		return nil, 0, err
	}

	return assignments, total, nil
}

func (r *Repository) UpdateAssignment(a *models.Assignment) (*models.Assignment, error) {
	if err := r.db.Save(a).Error; err != nil {
		return nil, err
	}
	return a, nil
}

func (r *Repository) DeleteAssignment(id uint) error {
	return r.db.Delete(&models.Assignment{}, id).Error
}

// --- Submissions ---

func (r *Repository) CreateSubmission(s *models.AssignmentSubmission) (*models.AssignmentSubmission, error) {
	if err := r.db.Create(s).Error; err != nil {
		return nil, err
	}
	return s, nil
}

func (r *Repository) GetSubmissionByID(id uint) (*models.AssignmentSubmission, error) {
	var s models.AssignmentSubmission
	if err := r.db.Preload("Student").Preload("Assignment.Teacher").First(&s, id).Error; err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *Repository) GetStudentSubmissions(studentID uint, page, limit int) ([]models.AssignmentSubmission, int64, error) {
	var submissions []models.AssignmentSubmission
	var total int64

	query := r.db.Where("student_id = ?", studentID)

	if err := query.Model(&models.AssignmentSubmission{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := query.Preload("Assignment").
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&submissions).Error; err != nil {
		return nil, 0, err
	}

	return submissions, total, nil
}

func (r *Repository) GetAssignmentSubmissions(assignmentID uint, page, limit int) ([]models.AssignmentSubmission, int64, error) {
	var submissions []models.AssignmentSubmission
	var total int64

	query := r.db.Where("assignment_id = ?", assignmentID)

	if err := query.Model(&models.AssignmentSubmission{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := query.Preload("Student").
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&submissions).Error; err != nil {
		return nil, 0, err
	}

	return submissions, total, nil
}

func (r *Repository) UpdateSubmissionAfterEval(s *models.AssignmentSubmission) (*models.AssignmentSubmission, error) {
	if err := r.db.Save(s).Error; err != nil {
		return nil, err
	}
	return s, nil
}

func (r *Repository) GetExistingSubmission(studentID, assignmentID uint) (*models.AssignmentSubmission, error) {
	var s models.AssignmentSubmission
	if err := r.db.Where("student_id = ? AND assignment_id = ?", studentID, assignmentID).
		First(&s).Error; err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *Repository) GetStudentAssignmentStats(studentID uint) (totalSubmitted int64, avgScore float64, err error) {
	// Total submitted
	if err := r.db.Model(&models.AssignmentSubmission{}).
		Where("student_id = ? AND is_evaluated = ?", studentID, true).
		Count(&totalSubmitted).Error; err != nil {
		return 0, 0, err
	}

	// Average score
	var result struct {
		Avg float64
	}
	if err := r.db.Model(&models.AssignmentSubmission{}).
		Select("AVG(score) as avg").
		Where("student_id = ? AND is_evaluated = ?", studentID, true).
		Scan(&result).Error; err != nil {
		return totalSubmitted, 0, err
	}

	return totalSubmitted, result.Avg, nil
}
