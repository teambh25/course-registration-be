package service

import (
	"course-reg/internal/app/models"
	"course-reg/internal/pkg/constant"
)

type AdminServiceInterface interface {
	RegisterStudents([]models.Student) error
	ResetStudents() error
	RegisterCourses([]models.Course) error
	ResetCourses() error
	CreateCourse(*models.Course) (uint, error)
	DeleteCourse(uint) error

	GetRegistrationState() bool
	StartRegistration() error
	PauseRegistration() error
	GetRegistrationPeriod() (string, string)
	SetRegistrationPeriod(string, string) error
}

type AuthServiceInterface interface {
	Check(username string, password string) (constant.UserRole, uint, error)
}

type CourseRegServiceInterface interface {
	Enroll(studentID, courseID uint) (bool, string)
	// CancelEnrollment(studentID, courseID uint) (success bool, message string, allSeats map[uint]int)
}
