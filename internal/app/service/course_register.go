package service

import (
	"course-reg/internal/app/models"
	"course-reg/internal/app/repository"
	"course-reg/internal/app/worker"
	"course-reg/internal/pkg/setting"
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

func (s *CourseRegService) Enroll(studentID, courseID uint) (bool, string, map[uint]int, int) {
	// Check if within registration period
	withinPeriod, err := setting.IsWithinRegistrationPeriod()
	if err != nil {
		log.Println("failed to check registration period:", err.Error())
		return false, "수강신청 기간 확인 실패", nil, 0
	}
	if !withinPeriod {
		return false, "수강신청 기간이 아닙니다", nil, 0
	}

	resp := s.enrollmentWorker.Enroll(studentID, courseID)
	return resp.Success, resp.Message, resp.AllRemainingSeats, resp.WaitlistPosition
}

func (s *CourseRegService) CancelEnrollment(studentID, courseID uint) (bool, string, map[uint]int) {
	// Check if within registration period
	withinPeriod, err := setting.IsWithinRegistrationPeriod()
	if err != nil {
		log.Println("failed to check registration period:", err.Error())
		return false, "수강신청 기간 확인 실패", nil
	}
	if !withinPeriod {
		return false, "수강신청 기간이 아닙니다", nil
	}

	resp := s.enrollmentWorker.CancelEnrollment(studentID, courseID)
	return resp.Success, resp.Message, resp.AllRemainingSeats
}
