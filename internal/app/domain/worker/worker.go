package worker

import (
	"course-reg/internal/app/domain/cache"
	"course-reg/internal/app/repository"
	"sync"
)

type RequestType int

const (
	ENROLL RequestType = iota + 1
	CANCEL
	READ_ALL
	ADMIN_ENROLL
	ADMIN_CANCEL
)

// EnrollmentWorker handles enrollment operations with cache
type EnrollmentWorker struct {
	wg          sync.WaitGroup
	queueSize   int
	requestChan chan EnrollmentRequest
	cache       *cache.EnrollmentCache
	enrollRepo  repository.EnrollmentRepositoryInterface
}

func NewEnrollmentWorker(queueSize int, enrollRepo repository.EnrollmentRepositoryInterface) *EnrollmentWorker {
	return &EnrollmentWorker{
		queueSize:  queueSize,
		enrollRepo: enrollRepo,
	}
}
