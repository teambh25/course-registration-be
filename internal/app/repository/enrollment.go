package repository

import (
	"course-reg/internal/app/models"
	"fmt"

	"gorm.io/gorm"
)

type EnrollmentRepository struct {
	db *gorm.DB
}

func NewEnrollmentRepository(db *gorm.DB) *EnrollmentRepository {
	return &EnrollmentRepository{db: db}
}

func (r *EnrollmentRepository) InsertEnrollment(enrollment *models.Enrollment) error {
	result := r.db.Create(enrollment)
	if result.Error != nil {
		return fmt.Errorf("create failed: %w", result.Error)
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
