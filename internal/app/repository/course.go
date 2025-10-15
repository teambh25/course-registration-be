package repository

import (
	"course-reg/internal/app/models"
	"fmt"

	"gorm.io/gorm"
)

type CourseRepository struct {
	db *gorm.DB
}

func NewCourseRepository(db *gorm.DB) *CourseRepository {
	return &CourseRepository{db: db}
}

func (r *CourseRepository) CreateCourse(course *models.Course) error {
	result := r.db.Create(course)
	if result.Error != nil {
		return fmt.Errorf("create failed: %w", result.Error)
	}
	return nil
}

func (r *CourseRepository) DeleteCourse(courseID uint) error {
	result := r.db.Delete(&models.Course{}, courseID)
	if result.Error != nil {
		return fmt.Errorf("delete failed: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("course not found") // todo: 커스텀 예외
	}
	return nil
}
