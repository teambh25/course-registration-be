package repository

import (
	"course-reg/internal/app/models"
	"fmt"

	"gorm.io/gorm"
)

type StudentRepository struct {
	db *gorm.DB
}

func NewStudentRepository(db *gorm.DB) *StudentRepository {
	return &StudentRepository{db: db}
}

func (r *StudentRepository) GetPassword(username string) (string, error) {
	var student *models.Student

	result := r.db.Where("phone_number = ?", username).Take(&student)

	// 중복으로 인한 실패인지 랩핑 필요
	return student.BirthDate, result.Error
}

func (r *StudentRepository) BulkInsertStudents(students []models.Student) error {
	tx := r.db.Begin()
	if err := tx.CreateInBatches(students, 100).Error; err != nil {
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
		return nil, fmt.Errorf("select failed: %w", result.Error)
	}
	return students, nil
}
