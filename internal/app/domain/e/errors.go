package e

import (
	"errors"
)

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrInvalidInput     = errors.New("invalid input data")
	ErrDuplicateStudent = errors.New("duplicate student")

	// Worker
	ErrWorkerInternal = errors.New("internal worker error")
	ErrWorkerTimeout  = errors.New("worker request timeout")

	// Enrollment
	ErrCourseNotFound            = errors.New("course not found")
	ErrStudentNotFound           = errors.New("student not found")
	ErrTimeConflict              = errors.New("time conflict with enrolled course")
	ErrAlreadyEnrolled           = errors.New("already enrolled in this course")
	ErrCourseFull                = errors.New("course is full")
	ErrEnrollmentDBFailed        = errors.New("failed to save enrollment")
	ErrCacheSyncFailed           = errors.New("cache sync failed")
	ErrInvalidRegistrationPeriod = errors.New("not within registration period")

	// DB constraint errors
	ErrDBDuplicateEnrollment = errors.New("db: duplicate enrollment")
	ErrDBPositionTaken       = errors.New("db: position already taken")
)
