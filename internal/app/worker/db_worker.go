package worker

import (
	"course-reg/internal/app/models"
	"course-reg/internal/app/repository"
	"log"
)

// DBOperationType defines types of DB operations
type DBOperationType int

const (
	DB_SAVE DBOperationType = iota
	DB_DELETE
	DB_PROMOTE // delete waitlist + save enrollment
)

// DBOperation represents a DB operation request
type DBOperation struct {
	Type       DBOperationType
	StudentID  uint
	CourseID   uint
	IsWaitlist bool
	Position   int
}

// DBWorker handles all DB operations sequentially
type DBWorker struct {
	dbChan     chan DBOperation
	enrollRepo repository.EnrollmentRepositoryInterface
}

func NewDBWorker(enrollRepo repository.EnrollmentRepositoryInterface) *DBWorker {
	return &DBWorker{
		dbChan:     make(chan DBOperation, 1000),
		enrollRepo: enrollRepo,
	}
}

// Start starts the DB worker goroutine
func (w *DBWorker) Start() {
	go w.worker()
}

// worker processes all DB operations sequentially (no concurrent writes to SQLite)
func (w *DBWorker) worker() {
	for op := range w.dbChan {
		switch op.Type {
		case DB_SAVE:
			enrollment := &models.Enrollment{
				StudentID:        op.StudentID,
				CourseID:         op.CourseID,
				IsWaitlist:       op.IsWaitlist,
				WaitlistPosition: op.Position,
			}
			if err := w.enrollRepo.SaveEnrollment(enrollment); err != nil {
				log.Printf("Failed to save enrollment to DB: %v", err)
			}

		case DB_DELETE:
			if err := w.enrollRepo.DeleteEnrollment(op.StudentID, op.CourseID); err != nil {
				log.Printf("Failed to delete enrollment from DB: %v", err)
			}

		case DB_PROMOTE:
			// Delete waitlist entry
			if err := w.enrollRepo.DeleteEnrollment(op.StudentID, op.CourseID); err != nil {
				log.Printf("Failed to delete waitlist entry: %v", err)
			}
			// Save as enrollment
			enrollment := &models.Enrollment{
				StudentID:  op.StudentID,
				CourseID:   op.CourseID,
				IsWaitlist: false,
			}
			if err := w.enrollRepo.SaveEnrollment(enrollment); err != nil {
				log.Printf("Failed to save promoted enrollment: %v", err)
			}
		}
	}
}

// SendOperation sends DB operation to worker (non-blocking)
func (w *DBWorker) SendOperation(op DBOperation) {
	w.dbChan <- op
}
