package service

import (
	"course-reg/internal/app/domain/worker"
	"course-reg/internal/app/models"
	"course-reg/internal/app/repository"
	"course-reg/internal/pkg/constant"
	"log"
)

type CourseRegService struct {
	courseRepo       repository.CourseRepositoryInterface
	enrollRepo       repository.EnrollmentRepositoryInterface
	enrollmentWorker *worker.EnrollmentWorker
}

func NewCourseRegService(
	c repository.CourseRepositoryInterface,
	e repository.EnrollmentRepositoryInterface,
	w *worker.EnrollmentWorker,
) *CourseRegService {
	return &CourseRegService{
		courseRepo:       c,
		enrollRepo:       e,
		enrollmentWorker: w,
	}
}

func (s *CourseRegService) GetAllCourses() ([]models.Course, error) {
	courses, err := s.courseRepo.FetchAllCourses()
	if err != nil {
		log.Println("fetch all courses failed:", err.Error())
	}
	return courses, err
}

func (s *CourseRegService) Enroll(studentID, courseID uint) (bool, string, map[uint]constant.CourseStatus, int) {
	resp := s.enrollmentWorker.Enroll(studentID, courseID)
	return resp.Success, resp.Message, resp.CourseStatuses, resp.WaitlistPosition
}

// func (s *CourseRegService) CancelEnrollment(studentID, courseID uint) (bool, string, map[uint]int) {
// 	resp := s.enrollmentWorker.CancelEnrollment(studentID, courseID)
// 	return resp.Success, resp.Message, resp.CourseStatuses
// }
