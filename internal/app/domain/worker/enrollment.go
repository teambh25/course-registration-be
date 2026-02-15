package worker

import (
	"context"
	"course-reg/internal/app/domain/cache"
	"course-reg/internal/app/domain/constants"
	"course-reg/internal/app/domain/e"
	"course-reg/internal/app/models"
	"errors"
	"fmt"
	"log"
)

// EnrollmentRequest represents an enrollment request
type EnrollmentRequest struct {
	Type      RequestType
	StudentID uint
	CourseID  uint
	Response  chan error
}

func (w *EnrollmentWorker) Start(students []models.Student, courses []models.Course, enrollments []models.Enrollment) error {
	if w.requestChan != nil {
		return errors.New("worker already running")
	}

	enrollmentCache, err := cache.NewEnrollmentCacheWithData(students, courses, enrollments)
	if err != nil {
		return err
	}

	w.requestChan = make(chan EnrollmentRequest, w.queueSize)
	w.cache = enrollmentCache

	w.wg.Add(1)
	go func() {
		defer w.wg.Done()
		w.worker()
	}()

	return nil
}

func (w *EnrollmentWorker) Stop() {
	close(w.requestChan)
	w.wg.Wait()
	w.requestChan = nil
}

func (w *EnrollmentWorker) worker() {
	for req := range w.requestChan {
		var err error

		// Use defer inside an anonymous function to ensure resources are released promptly
		// and to prevent potential resource leaks within the loop scope.
		func() {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("[fatal] worker panic: type=%d studentID=%d courseID=%d err=%v", req.Type, req.StudentID, req.CourseID, r)
					err = e.ErrWorkerInternal
				}
			}()

			ctx, cancel := context.WithTimeout(context.Background(), workerTimeout)
			defer cancel()

			switch req.Type {
			case ENROLL:
				err = w.processEnroll(ctx, req)
			}

			if err != nil && ctx.Err() != nil {
				err = fmt.Errorf("%w: %v", e.ErrWorkerTimeout, err)
			}
		}()

		req.Response <- err
	}
}

func (w *EnrollmentWorker) Enroll(studentID, courseID uint) error {
	req := EnrollmentRequest{
		Type:      ENROLL,
		StudentID: studentID,
		CourseID:  courseID,
		Response:  make(chan error, 1),
	}

	w.requestChan <- req
	return <-req.Response
}

// processEnroll handles enrollment logic
func (w *EnrollmentWorker) processEnroll(ctx context.Context, req EnrollmentRequest) error {
	studentID := req.StudentID
	courseID := req.CourseID

	if !w.cache.StudentExists(studentID) {
		log.Printf("[fatal] invalide student ID : type=%d studentID=%d courseID=%d err=%v", req.Type, req.StudentID, req.CourseID)
		return e.ErrStudentNotFound
	}

	if !w.cache.CourseExists(courseID) {
		return e.ErrCourseNotFound
	}

	if w.cache.IsStudentEnrolled(studentID, courseID) {
		return e.ErrAlreadyEnrolled
	}

	if w.cache.HasTimeConflict(studentID, courseID) {
		return e.ErrTimeConflict
	}

	pos, ok := w.cache.GetAvailablePos(courseID)
	if !ok {
		return e.ErrCourseFull
	}

	if err := w.enrollRepo.InsertEnrollment(ctx, &models.Enrollment{StudentID: studentID, CourseID: courseID, Position: pos}); err != nil {
		return fmt.Errorf("%w: %v", e.ErrEnrollmentDBFailed, err)
	}
	w.cache.EnrollStudent(studentID, courseID)

	return nil
}

// todo : 다른 파일로 분리?
func (w *EnrollmentWorker) GetAllCourseStatus() map[uint]constants.CourseStatus {
	status := make(map[uint]constants.CourseStatus)
	for courseID, info := range w.cache.GetAllCourseCountInfo() {
		if info.EnrolledCount < info.Capacity {
			status[courseID] = constants.CourseAvailable
		} else if info.WaitingCount < info.Capacity {
			status[courseID] = constants.CourseWaitlist
		} else {
			status[courseID] = constants.CourseFull
		}
	}
	return status
}

func (w *EnrollmentWorker) processAddWaitList() {
	// isWaitlistFull, err := w.cache.IsWaitlistFull(courseID)
	// if err != nil {
	// 	return EnrollmentResponse{
	// 		Success:        false,
	// 		Message:        "강의 정보를 찾을 수 없습니다",
	// 	}
	// }

	// if !isWaitlistFull {
	// 	position, err := w.cache.AddToWaitlist(studentID, courseID)
	// 	if err != nil {
	// 		return EnrollmentResponse{
	// 			Success:        false,
	// 			Message:        "존재하지 않는 학생입니다",
	// 		}
	// 	}
	// 	return EnrollmentResponse{
	// 		Success:          false,
	// 		Message:          "대기열에 등록되었습니다",
	// 		WaitlistPosition: position,
	// 	}
	// }

	// // Both course and waitlist are full
	// return EnrollmentResponse{
	// 	Success:        false,
	// 	Message:        "대기열이 마감되었습니다",
	// }

}
