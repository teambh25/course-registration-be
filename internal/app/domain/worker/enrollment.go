package worker

import (
	"course-reg/internal/app/domain/cache"
	"course-reg/internal/app/models"
	"course-reg/internal/app/repository"
)

// EnrollmentWorker handles enrollment operations with cache
type EnrollmentWorker struct {
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

// EnrollmentRequest represents an enrollment request
type EnrollmentRequest struct {
	Type      RequestType
	StudentID uint
	CourseID  uint
	Response  chan EnrollmentResponse
}

// EnrollmentResponse represents the response to an enrollment request
type EnrollmentResponse struct {
	Success bool
	Message string
}

func NewEnrollmentWorker(queueSize int) *EnrollmentWorker {
	return &EnrollmentWorker{
		requestChan: make(chan EnrollmentRequest, queueSize),
		cache:       cache.NewEnrollmentCache(),
	}
}

func (w *EnrollmentWorker) Start(students []models.Student, courses []models.Course, enrollments []models.Enrollment) {
	w.cache.LoadInitStudents(students)
	w.cache.LoadInitCourses(courses)
	w.cache.LoadEnrollments(enrollments)
	w.cache.BuildConflictGraph(courses)
	go w.worker()
}

func (w *EnrollmentWorker) worker() {
	for req := range w.requestChan {
		var response EnrollmentResponse

		switch req.Type {
		case ENROLL:
			response = w.processEnroll(req)
			// case READ_ALL:
			// 	response = w.processReadAll()
			// case CANCEL:
			// 	response = w.processCancel(req)
			// case ADMIN_ENROLL:
			// 	response = w.processAdminEnroll(req)
		}

		req.Response <- response
	}
}

func (w *EnrollmentWorker) Enroll(studentID, courseID uint) EnrollmentResponse {
	req := EnrollmentRequest{
		Type:      ENROLL,
		StudentID: studentID,
		CourseID:  courseID,
		Response:  make(chan EnrollmentResponse, 1),
	}

	w.requestChan <- req
	return <-req.Response
}

// processEnroll handles enrollment logic
func (w *EnrollmentWorker) processEnroll(req EnrollmentRequest) EnrollmentResponse {
	studentID := req.StudentID
	courseID := req.CourseID

	// Todo : validate handler쪽으로 옮기기
	if !w.cache.CourseExists(courseID) {
		return EnrollmentResponse{
			Success: false,
			Message: "존재하지 않는 강의입니다",
		}
	}

	// Todo : 500 에러 처리
	if !w.cache.StudentExists(studentID) {
		return EnrollmentResponse{
			Success: false,
			Message: "존재하지 않는 학생입니다",
		}
	}

	if w.cache.HasTimeConflict(studentID, courseID) {
		return EnrollmentResponse{
			Success: false,
			Message: "시간이 겹치는 강의가 있습니다",
		}
	}

	if w.cache.IsStudentEnrolled(studentID, courseID) {
		return EnrollmentResponse{
			Success: false,
			Message: "이미 신청한 강의입니다",
		}
	}

	pos, err := w.cache.GetPosIfNotFull(courseID)
	if err != nil {
		return EnrollmentResponse{
			Success: false,
			Message: "정원이 초과 되었습니다",
		}
	}

	w.enrollRepo.InsertEnrollment(&models.Enrollment{StudentID: studentID, CourseID: courseID, Position: pos})
	w.cache.EnrollStudent(studentID, courseID)

	return EnrollmentResponse{
		Success: true,
		Message: "수강신청 성공",
	}
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
