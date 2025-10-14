package repository

import (
	"course-reg/models"

	"gorm.io/gorm"
)

type StudentRepository struct {
	db *gorm.DB
}

func NewStudentRepository(db *gorm.DB) *StudentRepository {
	return &StudentRepository{db: db}
}

func (r *StudentRepository) InsertStudens(students []models.Student) error {
	tx := r.db.Begin()
	if err := tx.CreateInBatches(students, 100).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (r *StudentRepository) GetPassword(username string) (string, error) {
	var student *models.Student

	result := r.db.Where("phone_number = ?", username).Take(&student)
	return student.BirthDate, result.Error
}
