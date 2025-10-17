package repository

import "course-reg/internal/app/models"

type StudentRepositoryInterface interface {
	GetPassword(username string) (string, error)
	BulkInsertStudents(students []models.Student) error
	DeleteAllStudents() error
}

type CourseRepositoryInterface interface {
	BulkInsertCourses(courses []models.Course) error
	DeleteAllCourses() error
	CreateCourse(course *models.Course) error
	DeleteCourse(courseID uint) error
	FetchAllCourses() ([]models.Course, error)
}

type EnrollmentRepositoryInterface interface {
	SaveEnrollment(enrollment *models.Enrollment) error
	DeleteEnrollment(studentID uint, courseID uint) error
	LoadAllEnrollments() ([]models.Enrollment, error)
}
