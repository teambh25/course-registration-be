package service

import (
	"course-reg/internal/app/models"
	"course-reg/internal/app/repository"
	"log"
)

type CourseRegService struct {
	courseRepo repository.CourseRepositoryInterface
	enrollRepo repository.EnrollmentRepositoryInterface
}

func NewCourseRegService(
	c repository.CourseRepositoryInterface,
	e repository.EnrollmentRepositoryInterface,
) *CourseRegService {
	return &CourseRegService{courseRepo: c, enrollRepo: e}
}

func (s *CourseRegService) GetAllCourses() ([]models.Course, error) {
	courses, err := s.courseRepo.FetchAllCourses()
	if err != nil {
		log.Println("fetch all courses failed:", err.Error())
	}
	return courses, err
}
