package e

import (
	"errors"
)

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrInvalidInput     = errors.New("invalid input data")
	ErrDuplicateStudent = errors.New("duplicate student")

	// for Course Registration
	ErrCourseNotFound            = errors.New("course not found")
	ErrStudentNotFound           = errors.New("student not found")
	ErrTimeConflict              = errors.New("time conflict with enrolled course")
	ErrAlreadyEnrolled           = errors.New("already enrolled in this course")
	ErrCourseFull                = errors.New("course is full")
	ErrInvalidRegistrationPeriod = errors.New("not within registration period")
)
