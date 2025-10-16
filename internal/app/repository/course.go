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

func (r *CourseRepository) BulkInsertCourses(courses []models.Course) error {
	tx := r.db.Begin()
	if err := tx.CreateInBatches(courses, 100).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("create in batches failed: %w", err)
	}
	tx.Commit()
	return nil
}

func (r *CourseRepository) DeleteAllCourses() error {
	if err := r.db.Migrator().DropTable(&models.Course{}); err != nil {
		return fmt.Errorf("drop table failed: %w", err)
	}
	if err := r.db.AutoMigrate(&models.Course{}); err != nil {
		return fmt.Errorf("auto migrate failed: %w", err)
	}
	return nil
}

func (r *CourseRepository) CreateCourse(course *models.Course) error {
	result := r.db.Create(course)
	if result.Error != nil {
		return fmt.Errorf("create failed: %w", result.Error)
	}
	return nil
}

func (r *CourseRepository) FetchAllCourses() ([]models.Course, error) {
	var courses []models.Course
	result := r.db.Find(&courses)
	if result.Error != nil {
		return nil, fmt.Errorf("select failed: %w", result.Error)
	}
	return courses, nil
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
