package repository

import "gorm.io/gorm"

type EnrollmentRepository struct {
	db *gorm.DB
}

func NewEnrollmentRepositoryRepository(db *gorm.DB) *EnrollmentRepository {
	return &EnrollmentRepository{db: db}
}
