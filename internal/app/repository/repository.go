package repository

import "course-reg/internal/app/models"

type StudentRepositoryInterface interface {
	GetPassword(username string) (string, error)
	InsertStudents(students []models.Student) error
	DeleteAllStudents() error
}

type CourseRepositoryInterface interface {
	CreateCourse(course *models.Course) error
	DeleteCourse(courseID uint) error
	FetchAllCourses() ([]models.Course, error)
}

type EnrollmentRepositoryInterface interface {
}
