package repository

import (
	"course-reg/internal/app/models"

	"gorm.io/gorm"
)

type EnrollmentRepository struct {
	db *gorm.DB
}

func NewEnrollmentRepository(db *gorm.DB) *EnrollmentRepository {
	return &EnrollmentRepository{db: db}
}

func (r *EnrollmentRepository) SaveEnrollment(enrollment *models.Enrollment) error {
	return r.db.Create(enrollment).Error
}

func (r *EnrollmentRepository) DeleteEnrollment(studentID uint, courseID uint) error {
	return r.db.Where("student_id = ? AND course_id = ?", studentID, courseID).Delete(&models.Enrollment{}).Error
}

func (r *EnrollmentRepository) LoadAllEnrollments() ([]models.Enrollment, error) {
	var enrollments []models.Enrollment
	err := r.db.Find(&enrollments).Error
	return enrollments, err
}
