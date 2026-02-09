package service

import (
	"course-reg/internal/app/domain/cache"
	"course-reg/internal/app/domain/worker"
	"course-reg/internal/app/repository"
	"course-reg/internal/pkg/constant"
)

type CourseRegService struct {
	courseRepo       repository.CourseRepositoryInterface
	enrollRepo       repository.EnrollmentRepositoryInterface
	enrollmentWorker *worker.EnrollmentWorker
	regState         *cache.RegistrationState
}

func NewCourseRegService(
	c repository.CourseRepositoryInterface,
	e repository.EnrollmentRepositoryInterface,
	w *worker.EnrollmentWorker,
	r *cache.RegistrationState,
) *CourseRegService {
	return &CourseRegService{
		courseRepo:       c,
		enrollRepo:       e,
		enrollmentWorker: w,
		regState:         r,
	}
}

func (s *CourseRegService) Enroll(studentID, courseID uint) worker.EnrollmentResult {
	var result worker.EnrollmentResult

	err := s.regState.RunIfEnabled(true, func() error {
		result = s.enrollmentWorker.Enroll(studentID, courseID)
		return nil
	})

	if err != nil {
		return worker.EnrollNotInPeriod
	}

	return result
}

func (s *CourseRegService) GetAllCourseStatus() (map[uint]constant.CourseStatus, error) {
	var result map[uint]constant.CourseStatus

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
