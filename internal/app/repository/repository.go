package repository

import "course-reg/internal/app/models"

type StudentRepositoryInterface interface {
	FetchPassword(username string) (uint, string, error)
	BulkInsertStudents(students []models.Student) error
	DeleteAllStudents() error
	FetchAllStudents() ([]models.Student, error)
}

type CourseRepositoryInterface interface {
	BulkInsertCourses(courses []models.Course) error
	DeleteAllCourses() error
	InsertCourse(course *models.Course) error
	DeleteCourse(courseID uint) error
	FetchAllCourses() ([]models.Course, error)
}

type EnrollmentRepositoryInterface interface {
	InsertEnrollment(enrollment *models.Enrollment) error
	DeleteEnrollment(studentID uint, courseID uint) error
	FetchAllEnrollments() ([]models.Enrollment, error)
}
