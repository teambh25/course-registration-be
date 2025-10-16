package service

import "course-reg/internal/app/models"

type AdminServiceInterface interface {
	RegisterStudents(students []models.Student) error
	ResetStudents() error
	RegisterCourses(courses []models.Course) error
	ResetCourses() error
	CreateCourse(course *models.Course) (uint, error)
	DeleteCourse(courseID uint) error
}

type CourseRegServiceInterface interface {
	GetAllCourses() ([]models.Course, error)
}
