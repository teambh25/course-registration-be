package service

import (
	"course-reg/internal/app/domain/cache"
	"course-reg/internal/app/domain/constants"
	"course-reg/internal/app/domain/worker"
	"course-reg/internal/app/repository"
)

type CourseRegService struct {
	courseRepo       repository.CourseRepositoryInterface
	enrollRepo       repository.EnrollmentRepositoryInterface
	enrollmentWorker *worker.Worker
	regState         *cache.RegistrationState
}

func NewCourseRegService(
	c repository.CourseRepositoryInterface,
	e repository.EnrollmentRepositoryInterface,
	w *worker.Worker,
	r *cache.RegistrationState,
) *CourseRegService {
	return &CourseRegService{
		courseRepo:       c,
		enrollRepo:       e,
		enrollmentWorker: w,
		regState:         r,
	}
}

func (s *CourseRegService) Enroll(studentID, courseID uint) error {
	return s.regState.RunIfEnabled(true, func() error {
		return s.enrollmentWorker.Enroll(studentID, courseID)
	})
}

func (s *CourseRegService) GetAllCourseStatus() (map[uint]constants.CourseStatus, error) {
	var result map[uint]constants.CourseStatus
	err := s.regState.RunIfEnabled(true, func() error {
		result = s.enrollmentWorker.GetAllCourseStatus()
		return nil
	})
	return result, err
}

// func (s *CourseRegService) CancelEnrollment(studentID, courseID uint) (bool, string, map[uint]int) {
// 	resp := s.enrollmentWorker.CancelEnrollment(studentID, courseID)
// 	return resp.Success, resp.Message, resp.CourseStatuses
// }
