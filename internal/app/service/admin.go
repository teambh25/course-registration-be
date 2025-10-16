package service

import (
	"course-reg/internal/app/models"
	"course-reg/internal/app/repository"
	"log"
)

type AdminService struct {
	studentRepo repository.StudentRepositoryInterface
	courseRepo  repository.CourseRepositoryInterface
	enrollRepo  repository.EnrollmentRepositoryInterface
}

func NewAdminService(
	s repository.StudentRepositoryInterface,
	c repository.CourseRepositoryInterface,
	e repository.EnrollmentRepositoryInterface,
) *AdminService {
	return &AdminService{
		studentRepo: s,
		courseRepo:  c,
		enrollRepo:  e,
	}
}

func (s *AdminService) RegisterStudents(students []models.Student) error {
	err := s.studentRepo.InsertStudents(students)
	if err != nil {
		log.Println("register students failed:", err.Error())
	}
	return err
}

func (s *AdminService) ResetStudents() error {
	err := s.studentRepo.DeleteAllStudents()
	if err != nil {
		log.Println("reset students failed:", err.Error())
	}
	return err
}

func (s *AdminService) CreateCourse(course *models.Course) (uint, error) {
	err := s.courseRepo.CreateCourse(course)
	if err != nil {
		log.Println("create course failed:", err.Error())
		return 0, err
	}
	return course.ID, nil
}

func (s *AdminService) DeleteCourse(courseID uint) error {
	err := s.courseRepo.DeleteCourse(courseID)
	if err != nil {
		log.Println("delete course failed:", err.Error())
	}
	return err
}

// func (s *AdminService) GetEnrolledStudentsByCourse(courseID uint) ([]Student, error)

// func (s *AdminService) ForceEnrollStudent(courseID uint, studentID uint) error
// func (s *AdminService) CancelEnrollment(courseID uint, studentID uint) error
// func (s *AdminService) CheckDuplicateCourses(studentID uint, courseID uint) ([]Course, error)
