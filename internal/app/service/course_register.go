package service

import (
	"course-reg/internal/app/domain/cache"
	"course-reg/internal/app/domain/worker"
	"course-reg/internal/app/repository"
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

func (s *CourseRegService) Enroll(studentID, courseID uint) (bool, string) {
	var resp worker.EnrollmentResponse

	err := s.regState.RunIfEnabled(true, func() error {
		resp = s.enrollmentWorker.Enroll(studentID, courseID)
		return nil
	})

	if err != nil {
		return false, "수강신청 기간이 아닙니다"
	}

	return resp.Success, resp.Message
}

// func (s *CourseRegService) CancelEnrollment(studentID, courseID uint) (bool, string, map[uint]int) {
// 	resp := s.enrollmentWorker.CancelEnrollment(studentID, courseID)
// 	return resp.Success, resp.Message, resp.CourseStatuses
// }
