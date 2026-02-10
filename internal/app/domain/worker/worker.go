package worker

import (
	"course-reg/internal/app/domain/cache"
	"course-reg/internal/app/repository"
	"sync"
)

// Worker handles enrollment operations with cache
type Worker struct {
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
