package cache

import (
	"course-reg/internal/app/models"
	"sync/atomic"
)

// EnrollmentCache is a simple in-memory data structure
type EnrollmentCache struct {

	// Course data
	CourseCapacity map[uint]int  // courseID -> capacity
	ConflictGraph  ConflictGraph // courseID -> conflicting courseIDs

	// Enrollment data (atomic count-based)
	StudentCourses        map[uint]map[uint]struct{} // studentID -> set of enrolled courseIDs
	StudentWaitingCourses map[uint]map[uint]struct{} // studentID -> set of waiting courseIDs
	EnrolledCount         map[uint]*atomic.Int32     // courseID -> count of enrolled students (atomic)
	WaitingCount          map[uint]*atomic.Int32     // courseID -> count of waiting students (atomic)
}

func NewEnrollmentCacheWithData(students []models.Student, courses []models.Course, enrollments []models.Enrollment) (*EnrollmentCache, error) {
	cache := &EnrollmentCache{
		CourseCapacity:        make(map[uint]int),
		StudentCourses:        make(map[uint]map[uint]struct{}),
		StudentWaitingCourses: make(map[uint]map[uint]struct{}),
		EnrolledCount:         make(map[uint]*atomic.Int32),
		WaitingCount:          make(map[uint]*atomic.Int32),
	}
	cache.loadInitStudents(students)
	if err := cache.loadInitCourses(courses); err != nil {
		return nil, err
	}
	cache.loadEnrollments(enrollments)
	return cache, nil
}

func (cache *EnrollmentCache) loadInitStudents(students []models.Student) {
	for _, s := range students {
		cache.StudentCourses[s.ID] = make(map[uint]struct{})
		cache.StudentWaitingCourses[s.ID] = make(map[uint]struct{})
	}
}

func (cache *EnrollmentCache) loadInitCourses(courses []models.Course) error {
	for _, c := range courses {
		cache.CourseCapacity[c.ID] = c.Capacity
		cache.EnrolledCount[c.ID] = &atomic.Int32{}
		cache.WaitingCount[c.ID] = &atomic.Int32{}
	}

	conflictGraph, err := BuildConflictGraph(courses)
	if err != nil {
		return err
	}
	cache.ConflictGraph = conflictGraph
	return nil
}

// loadEnrollments loads existing enrollments into cache
// Must be called after LoadInitStudents and LoadInitCourses
func (cache *EnrollmentCache) loadEnrollments(enrollments []models.Enrollment) {
	// todo : 대기열 로직 다시 고민하기
	for _, e := range enrollments {
		if e.IsWaitlist {
			// Update to max(current, position + 1)
			currentMax := cache.WaitingCount[e.CourseID].Load()
			if int32(e.Position+1) > currentMax {
				cache.WaitingCount[e.CourseID].Store(int32(e.Position + 1))
			}
			cache.StudentWaitingCourses[e.StudentID][e.CourseID] = struct{}{}
		} else {
			// Update to max(current, position + 1)
			currentMax := cache.EnrolledCount[e.CourseID].Load()
			if int32(e.Position+1) > currentMax {
				cache.EnrolledCount[e.CourseID].Store(int32(e.Position + 1))
			}
			cache.StudentCourses[e.StudentID][e.CourseID] = struct{}{}
		}
	}
}

// StudentExists checks if a student exists in cache
func (cache *EnrollmentCache) StudentExists(studentID uint) bool {
	return cache.StudentCourses[studentID] != nil
}

// CourseExists checks if a course exists in cache
func (cache *EnrollmentCache) CourseExists(courseID uint) bool {
	_, exists := cache.CourseCapacity[courseID]
	return exists
}

type CourseCountInfo struct {
	Capacity      int
	EnrolledCount int
	WaitingCount  int
}

func (cache *EnrollmentCache) GetAllCourseCountInfo() map[uint]CourseCountInfo {
	info := make(map[uint]CourseCountInfo)
	for courseID, capacity := range cache.CourseCapacity {
		info[courseID] = CourseCountInfo{
			Capacity:      capacity,
			EnrolledCount: int(cache.EnrolledCount[courseID].Load()),
			WaitingCount:  int(cache.WaitingCount[courseID].Load()),
		}
	}
	return info
}

// HasTimeConflict checks if enrolling in a course would create a time conflict
// Assumes student existence is already validated
func (cache *EnrollmentCache) HasTimeConflict(studentID, courseID uint) bool {
	for enrolledCourse := range cache.StudentCourses[studentID] {
		if cache.ConflictGraph[courseID][enrolledCourse] {
			return true
		}
	}
	return false
}

// IsStudentEnrolled checks if a student is already enrolled in a course
// Assumes student existence is already validated
func (cache *EnrollmentCache) IsStudentEnrolled(studentID, courseID uint) bool {
	_, exists := cache.StudentCourses[studentID][courseID]
	return exists
}

func (cache *EnrollmentCache) GetAvailablePos(courseID uint) (int, bool) {
	capacity := cache.CourseCapacity[courseID]
	enrolledCount := int(cache.EnrolledCount[courseID].Load())
	if enrolledCount >= capacity {
		return 0, false
	}
	return enrolledCount, true
}

// EnrollStudent enrolls a student in a course
// Assumes student and course existence is already validated
func (cache *EnrollmentCache) EnrollStudent(studentID, courseID uint) {
	cache.EnrolledCount[courseID].Add(1)
	cache.StudentCourses[studentID][courseID] = struct{}{}
}

// // IsWaitlistFull checks if a course's waitlist has reached capacity
// // Assumes course existence is already validated
// func (cache *EnrollmentCache) IsWaitlistFull(courseID uint) bool {
// 	capacity := cache.CourseCapacity[courseID]
// 	waitingCount := int(cache.WaitingCount[courseID].Load())
// 	return waitingCount >= capacity
// }

// // AddToWaitlist adds a student to a course's waitlist and returns their position
// // Assumes student and course existence is already validated
// func (cache *EnrollmentCache) AddToWaitlist(studentID, courseID uint) int {
// 	newCount := cache.WaitingCount[courseID].Add(1)
// 	cache.StudentWaitingCourses[studentID][courseID] = struct{}{}
// 	return int(newCount)
// }
