package service

import (
	"course-reg/models"
	"course-reg/repository"
)

type AdminService struct {
	studentRepo *repository.StudentRepository
	courseRepo  *repository.CourseRepository
	enrollRepo  *repository.EnrollmentRepository
}

func NewAdminService(
	s *repository.StudentRepository,
	c *repository.CourseRepository,
	e *repository.EnrollmentRepository,
) *AdminService {
	return &AdminService{
		studentRepo: s,
		courseRepo:  c,
		enrollRepo:  e,
	}
}

func (s *AdminService) RegisterStudents(students []models.Student) error {
	err := s.studentRepo.InsertStudents(students)
	return err
}

func (s *AdminService) CreateCourse() {

}

func (s *AdminService) DeleteCourse() {

}

// func (s *AdminService) GetEnrolledStudentsByCourse(courseID uint) ([]Student, error)

// func (s *AdminService) ForceEnrollStudent(courseID uint, studentID uint) error
// func (s *AdminService) CancelEnrollment(courseID uint, studentID uint) error
// func (s *AdminService) CheckDuplicateCourses(studentID uint, courseID uint) ([]Course, error)
