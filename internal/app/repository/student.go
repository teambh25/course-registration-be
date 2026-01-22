package repository

import (
	"course-reg/internal/app/models"
	"fmt"

	"gorm.io/gorm"
)

const studentBatchSize = 100

type StudentRepository struct {
	db *gorm.DB
}

func NewStudentRepository(db *gorm.DB) *StudentRepository {
	return &StudentRepository{db: db}
}

func (r *StudentRepository) FetchPassword(username string) (uint, string, error) {
	var student *models.Student

	result := r.db.Where("phone_number = ?", username).Take(&student)
	if result.Error != nil {
		return 0, "", fmt.Errorf("fetch failed: %w", result.Error)
	}
	return student.ID, student.BirthDate, nil
}

func (r *StudentRepository) BulkInsertStudents(students []models.Student) error {
	tx := r.db.Begin()
	if err := tx.CreateInBatches(students, studentBatchSize).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("create in batches failed: %w", err)
	}
	tx.Commit()
	return nil
}

func (r *StudentRepository) DeleteAllStudents() error {
	if err := r.db.Migrator().DropTable(&models.Student{}); err != nil {
		return fmt.Errorf("drop table failed: %w", err)
	}
	if err := r.db.AutoMigrate(&models.Student{}); err != nil {
		return fmt.Errorf("auto migrate failed: %w", err)
	}
	return nil
}

func (r *StudentRepository) FetchAllStudents() ([]models.Student, error) {
	var students []models.Student
	result := r.db.Find(&students)
	if result.Error != nil {
		return nil, fmt.Errorf("find failed: %w", result.Error)
	}
	return students, nil
}
