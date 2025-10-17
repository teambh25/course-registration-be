package worker

import (
	"course-reg/internal/app/models"
	"course-reg/internal/pkg/cache"
	"time"
)

// RequestType defines types of enrollment requests
type RequestType int

const (
	ENROLL RequestType = iota
	CANCEL
	READ_ALL
	ADMIN_ENROLL
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
	Success           bool
	Message           string
	AllRemainingSeats map[uint]int
	WaitlistPosition  int
}

// EnrollmentWorker handles enrollment operations with cache
type EnrollmentWorker struct {
	requestChan chan EnrollmentRequest
	cache       *cache.EnrollmentCache
}

func NewEnrollmentWorker(queueSize int) *EnrollmentWorker {
	return &EnrollmentWorker{
		requestChan: make(chan EnrollmentRequest, queueSize),
		cache:       cache.NewEnrollmentCache(),
	}
}

func (w *EnrollmentWorker) LoadInitStudents(students []models.Student) {
	w.cache.LoadInitStudents(students)
}

func (w *EnrollmentWorker) LoadInitCourses(courses []models.Course) {
	w.cache.LoadInitCourses(courses)
}

func (w *EnrollmentWorker) AddCourse(course models.Course) {
	w.cache.AddCourse(course)
}

func (w *EnrollmentWorker) RemoveCourse(courseID uint) {
	w.cache.RemoveCourse(courseID)
}

func (w *EnrollmentWorker) ClearAllCourses() {
	w.cache.ClearAllCourses()
}

func (w *EnrollmentWorker) ClearAllStudents() {
	w.cache.ClearAllStudents()
}

func (w *EnrollmentWorker) Start(students []models.Student, courses []models.Course) {
	// todo: cache pre-load
	w.LoadInitStudents(students)
	w.LoadInitCourses(courses)
	go w.worker()
}

// worker processes all enrollment requests sequentially
func (w *EnrollmentWorker) worker() {
	for req := range w.requestChan {
		var response EnrollmentResponse

		switch req.Type {
		case ENROLL:
			response = w.processEnroll(req)
		case CANCEL:
			response = w.processCancel(req)
		case READ_ALL:
			response = w.processReadAll()
		case ADMIN_ENROLL:
			response = w.processAdminEnroll(req)
		}

		req.Response <- response
	}
}

// Enroll sends enrollment request to worker
func (w *EnrollmentWorker) Enroll(studentID, courseID uint) EnrollmentResponse {
	req := EnrollmentRequest{
		Type:      ENROLL,
		StudentID: studentID,
		CourseID:  courseID,
		Response:  make(chan EnrollmentResponse, 1),
	}

	w.requestChan <- req

	select {
	case resp := <-req.Response:
		return resp
	case <-time.After(5 * time.Second):
		return EnrollmentResponse{
			Success: false,
			Message: "요청 타임아웃",
		}
	}
}

// GetAllRemainingSeats sends read request to worker
func (w *EnrollmentWorker) GetAllRemainingSeats() map[uint]int {
	req := EnrollmentRequest{
		Type:     READ_ALL,
		Response: make(chan EnrollmentResponse, 1),
	}

	w.requestChan <- req
	resp := <-req.Response
	return resp.AllRemainingSeats
}

// CancelEnrollment sends cancel request to worker
func (w *EnrollmentWorker) CancelEnrollment(studentID, courseID uint) EnrollmentResponse {
	req := EnrollmentRequest{
		Type:      CANCEL,
		StudentID: studentID,
		CourseID:  courseID,
		Response:  make(chan EnrollmentResponse, 1),
	}

	w.requestChan <- req

	select {
	case resp := <-req.Response:
		return resp
	case <-time.After(5 * time.Second):
		return EnrollmentResponse{
			Success: false,
			Message: "요청 타임아웃",
		}
	}
}

// processEnroll handles enrollment logic
func (w *EnrollmentWorker) processEnroll(req EnrollmentRequest) EnrollmentResponse {
	studentID := req.StudentID
	courseID := req.CourseID

	course, exists := w.cache.Courses[courseID]
	if !exists {
		return EnrollmentResponse{
			Success:           false,
			Message:           "존재하지 않는 강의입니다",
			AllRemainingSeats: w.cache.GetAllRemainingSeats(),
		}
	}

	// Check if already enrolled
	if w.cache.StudentCourses[studentID] != nil && w.cache.StudentCourses[studentID][courseID] {
		return EnrollmentResponse{
			Success:           false,
			Message:           "이미 신청한 강의입니다",
			AllRemainingSeats: w.cache.GetAllRemainingSeats(),
		}
	}

	// Check time conflicts
	if w.cache.StudentCourses[studentID] != nil {
		for enrolledCourse := range w.cache.StudentCourses[studentID] {
			if w.cache.ConflictGraph[courseID][enrolledCourse] {
				return EnrollmentResponse{
					Success:           false,
					Message:           "시간이 겹치는 강의가 있습니다",
					AllRemainingSeats: w.cache.GetAllRemainingSeats(),
				}
			}
		}
	}

	// Check capacity
	enrolled := w.cache.EnrolledStudents[courseID]
	if len(enrolled) < course.Capacity {
		// Success: enroll student
		w.cache.EnrolledStudents[courseID] = append(enrolled, studentID)
		if w.cache.StudentCourses[studentID] == nil {
			w.cache.StudentCourses[studentID] = make(map[uint]bool)
		}
		w.cache.StudentCourses[studentID][courseID] = true

		return EnrollmentResponse{
			Success:           true,
			Message:           "수강신청 성공",
			AllRemainingSeats: w.cache.GetAllRemainingSeats(),
		}
	}

	// Add to waitlist
	waiting := w.cache.WaitingStudents[courseID]
	if len(waiting) < course.Capacity {
		w.cache.WaitingStudents[courseID] = append(waiting, studentID)
		position := len(w.cache.WaitingStudents[courseID])

		return EnrollmentResponse{
			Success:           false,
			Message:           "대기열에 등록되었습니다",
			AllRemainingSeats: w.cache.GetAllRemainingSeats(),
			WaitlistPosition:  position,
		}
	}

	return EnrollmentResponse{
		Success:           false,
		Message:           "대기열이 마감되었습니다",
		AllRemainingSeats: w.cache.GetAllRemainingSeats(),
	}
}

// processCancel handles cancel logic with automatic waitlist processing
func (w *EnrollmentWorker) processCancel(req EnrollmentRequest) EnrollmentResponse {
	studentID := req.StudentID
	courseID := req.CourseID

	// Check if enrolled
	if w.cache.StudentCourses[studentID] == nil || !w.cache.StudentCourses[studentID][courseID] {
		return EnrollmentResponse{
			Success:           false,
			Message:           "등록하지 않은 강의입니다",
			AllRemainingSeats: w.cache.GetAllRemainingSeats(),
		}
	}

	// Remove from enrolled list
	enrolled := w.cache.EnrolledStudents[courseID]
	for i, sid := range enrolled {
		if sid == studentID {
			w.cache.EnrolledStudents[courseID] = append(enrolled[:i], enrolled[i+1:]...)
			break
		}
	}
	delete(w.cache.StudentCourses[studentID], courseID)

	// Process waitlist
	waiting := w.cache.WaitingStudents[courseID]
	for i, waitingStudentID := range waiting {
		// Check time conflicts for waitlisted student
		hasConflict := false
		if w.cache.StudentCourses[waitingStudentID] != nil {
			for enrolledCourse := range w.cache.StudentCourses[waitingStudentID] {
				if w.cache.ConflictGraph[courseID][enrolledCourse] {
					hasConflict = true
					break
				}
			}
		}

		if !hasConflict {
			// Enroll waitlisted student
			w.cache.WaitingStudents[courseID] = append(waiting[:i], waiting[i+1:]...)
			w.cache.EnrolledStudents[courseID] = append(w.cache.EnrolledStudents[courseID], waitingStudentID)
			if w.cache.StudentCourses[waitingStudentID] == nil {
				w.cache.StudentCourses[waitingStudentID] = make(map[uint]bool)
			}
			w.cache.StudentCourses[waitingStudentID][courseID] = true

			return EnrollmentResponse{
				Success:           true,
				Message:           "취소 완료. 대기자가 자동 등록되었습니다",
				AllRemainingSeats: w.cache.GetAllRemainingSeats(),
			}
		}
	}

	return EnrollmentResponse{
		Success:           true,
		Message:           "취소 완료",
		AllRemainingSeats: w.cache.GetAllRemainingSeats(),
	}
}

// processReadAll returns all remaining seats
func (w *EnrollmentWorker) processReadAll() EnrollmentResponse {
	return EnrollmentResponse{
		Success:           true,
		AllRemainingSeats: w.cache.GetAllRemainingSeats(),
	}
}

// processAdminEnroll handles admin force enrollment
func (w *EnrollmentWorker) processAdminEnroll(req EnrollmentRequest) EnrollmentResponse {
	studentID := req.StudentID
	courseID := req.CourseID

	// Force enroll (ignore capacity)
	w.cache.EnrolledStudents[courseID] = append(w.cache.EnrolledStudents[courseID], studentID)
	if w.cache.StudentCourses[studentID] == nil {
		w.cache.StudentCourses[studentID] = make(map[uint]bool)
	}
	w.cache.StudentCourses[studentID][courseID] = true

	return EnrollmentResponse{
		Success:           true,
		Message:           "관리자 강제 등록 완료",
		AllRemainingSeats: w.cache.GetAllRemainingSeats(),
	}
}
