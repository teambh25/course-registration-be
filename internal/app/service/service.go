package service

import (
	"course-reg/internal/app/models"
	"course-reg/internal/pkg/constant"
)

type AdminServiceInterface interface {
	RegisterStudents(students []models.Student) error
	ResetStudents() error
	RegisterCourses(courses []models.Course) error
	ResetCourses() error
	CreateCourse(course *models.Course) (uint, error)
	DeleteCourse(courseID uint) error
	SetRegistrationPeriod(startTime, endTime string) error
	GetRegistrationPeriod() (startTime, endTime string, err error)
}

type CourseRegServiceInterface interface {
	GetAllCourses() ([]models.Course, error)
	Enroll(studentID, courseID uint) (success bool, message string, allSeats map[uint]constant.CourseStatus, waitlistPos int)
	// CancelEnrollment(studentID, courseID uint) (success bool, message string, allSeats map[uint]int)
}
