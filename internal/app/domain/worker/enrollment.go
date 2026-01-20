package worker

import (
	"course-reg/internal/app/domain/cache"
	"course-reg/internal/app/models"
	"course-reg/internal/app/repository"
	"errors"
	"sync"
)

// EnrollmentWorker handles enrollment operations with cache
type EnrollmentWorker struct {
	wg          sync.WaitGroup
	queueSize   int
	requestChan chan EnrollmentRequest
	cache       *cache.EnrollmentCache
	enrollRepo  repository.EnrollmentRepositoryInterface
}

type RequestType int

const (
	ENROLL RequestType = iota + 1
	CANCEL
	READ_ALL
	ADMIN_ENROLL
	ADMIN_CANCEL
)

// EnrollmentResult represents the result of an enrollment operation
type EnrollmentResult int

const (
	EnrollSuccess EnrollmentResult = iota
	EnrollCourseNotFound
	EnrollStudentNotFound
	EnrollTimeConflict
	EnrollAlreadyEnrolled
	EnrollCourseFull
	EnrollNotInPeriod
)

var enrollResultMessages = map[EnrollmentResult]string{
	EnrollSuccess:         "수강신청 성공",
	EnrollCourseNotFound:  "존재하지 않는 강의입니다",
	EnrollStudentNotFound: "존재하지 않는 학생입니다",
	EnrollTimeConflict:    "시간이 겹치는 강의가 있습니다",
	EnrollAlreadyEnrolled: "이미 신청한 강의입니다",
	EnrollCourseFull:      "정원이 초과 되었습니다",
	EnrollNotInPeriod:     "수강신청 기간이 아닙니다",
}

func (r EnrollmentResult) String() string {
	if msg, ok := enrollResultMessages[r]; ok {
		return msg
	}
	return "알 수 없는 오류"
}

// EnrollmentRequest represents an enrollment request
type EnrollmentRequest struct {
	Type      RequestType
	StudentID uint
	CourseID  uint
	Response  chan EnrollmentResult
}

func NewEnrollmentWorker(queueSize int, enrollRepo repository.EnrollmentRepositoryInterface) *EnrollmentWorker {
	return &EnrollmentWorker{
		queueSize:  queueSize,
		enrollRepo: enrollRepo,
	}
}

func (w *EnrollmentWorker) Start(students []models.Student, courses []models.Course, enrollments []models.Enrollment) error {
	if w.requestChan != nil {
		return errors.New("worker already running")
	}

	w.requestChan = make(chan EnrollmentRequest, w.queueSize)
	w.cache = cache.NewEnrollmentCacheWithData(students, courses, enrollments)

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
		var result EnrollmentResult

		switch req.Type {
		case ENROLL:
			result = w.processEnroll(req)
			// case READ_ALL:
			// 	result = w.processReadAll()
			// case CANCEL:
			// 	result = w.processCancel(req)
			// case ADMIN_ENROLL:
			// 	result = w.processAdminEnroll(req)
		}

		req.Response <- result
	}
}

func (w *EnrollmentWorker) Enroll(studentID, courseID uint) EnrollmentResult {
	req := EnrollmentRequest{
		Type:      ENROLL,
		StudentID: studentID,
		CourseID:  courseID,
		Response:  make(chan EnrollmentResult, 1),
	}

	w.requestChan <- req
	return <-req.Response
}

// processEnroll handles enrollment logic
func (w *EnrollmentWorker) processEnroll(req EnrollmentRequest) EnrollmentResult {
	studentID := req.StudentID
	courseID := req.CourseID

	// Todo : validate handler쪽으로 옮기기
	if !w.cache.CourseExists(courseID) {
		return EnrollCourseNotFound
	}

	// Todo : 500 에러 처리
	if !w.cache.StudentExists(studentID) {
		return EnrollStudentNotFound
	}

	if w.cache.HasTimeConflict(studentID, courseID) {
		return EnrollTimeConflict
	}

	if w.cache.IsStudentEnrolled(studentID, courseID) {
		return EnrollAlreadyEnrolled
	}

	pos, err := w.cache.GetPosIfNotFull(courseID)
	if err != nil {
		return EnrollCourseFull
	}

	w.enrollRepo.InsertEnrollment(&models.Enrollment{StudentID: studentID, CourseID: courseID, Position: pos})
	w.cache.EnrollStudent(studentID, courseID)

	return EnrollSuccess
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
