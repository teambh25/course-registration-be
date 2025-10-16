package service

import "course-reg/internal/app/models"

type AdminServiceInterface interface {
	RegisterStudents(students []models.Student) error
	ResetStudents() error
	CreateCourse(course *models.Course) (uint, error)
	DeleteCourse(courseID uint) error
}
