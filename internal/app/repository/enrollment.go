package repository

import (
	"context"
	"course-reg/internal/app/domain/e"
	"course-reg/internal/app/models"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

const (
	pgUniqueViolation        = "23505"
	constraintStudentCourse  = "idx_student_course"
	constraintCoursePosition = "idx_course_position"
)

type EnrollmentRepository struct {
	db *gorm.DB
}

func NewEnrollmentRepository(db *gorm.DB) *EnrollmentRepository {
	return &EnrollmentRepository{db: db}
}

func (r *EnrollmentRepository) InsertEnrollment(ctx context.Context, enrollment *models.Enrollment) error {
	result := r.db.WithContext(ctx).Create(enrollment)
	if result.Error != nil {
		var pgErr *pgconn.PgError
		if errors.As(result.Error, &pgErr) && pgErr.Code == pgUniqueViolation {
			switch pgErr.ConstraintName {
			case constraintStudentCourse:
				return e.ErrDBDuplicateEnrollment
			case constraintCoursePosition:
				return e.ErrDBPositionTaken
			}
		}
		return fmt.Errorf("create failed: %w", result.Error)
	}
	return nil
}

func (r *EnrollmentRepository) GetMaxPosition(ctx context.Context, courseID uint) (int, error) {
	var maxPos *int
	err := r.db.WithContext(ctx).
		Model(&models.Enrollment{}).
		Where("course_id = ?", courseID).
		Select("MAX(position)").
		Scan(&maxPos).Error
	if err != nil {
		return 0, fmt.Errorf("get max position failed: %w", err)
	}
	if maxPos == nil {
		return -1, nil
	}
	return *maxPos, nil
}

func (r *EnrollmentRepository) BatchInsertEnrollments(enrollments []models.Enrollment) error {
	result := r.db.Create(&enrollments)
	if result.Error != nil {
		return fmt.Errorf("batch create failed: %w", result.Error)
	}
	return nil
}

func (r *EnrollmentRepository) DeleteEnrollment(studentID uint, courseID uint) error {
	result := r.db.Where("student_id = ? AND course_id = ?", studentID, courseID).Delete(&models.Enrollment{})
	if result.Error != nil {
		return fmt.Errorf("delete failed: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("enrollment not found") // todo: 커스텀 예외
	}
	return nil
}

func (r *EnrollmentRepository) DeleteAllEnrollments() error {
	if err := r.db.Migrator().DropTable(&models.Enrollment{}); err != nil {
		return fmt.Errorf("drop table failed: %w", err)
	}
	if err := r.db.AutoMigrate(&models.Enrollment{}); err != nil {
		return fmt.Errorf("auto migrate failed: %w", err)
	}
	return nil
}

func (r *EnrollmentRepository) FetchAllEnrollments() ([]models.Enrollment, error) {
	var enrollments []models.Enrollment
	err := r.db.Find(&enrollments).Error
	return enrollments, err
}
