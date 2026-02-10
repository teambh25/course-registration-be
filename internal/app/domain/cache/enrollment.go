package cache

import (
	"course-reg/internal/app/models"
	"errors"
	"fmt"
	"sync/atomic"
)

// EnrollmentCache is a simple in-memory data structure
type EnrollmentCache struct {

	// Course data
	CourseCapacity map[uint]int           // courseID -> capacity
	ConflictGraph  map[uint]map[uint]bool // courseID -> conflicting courseIDs

	// Enrollment data (atomic count-based)
	StudentCourses        map[uint]map[uint]struct{} // studentID -> set of enrolled courseIDs
	StudentWaitingCourses map[uint]map[uint]struct{} // studentID -> set of waiting courseIDs
	EnrolledCount         map[uint]*atomic.Int32     // courseID -> count of enrolled students (atomic)
	WaitingCount          map[uint]*atomic.Int32     // courseID -> count of waiting students (atomic)
}

func NewEnrollmentCacheWithData(students []models.Student, courses []models.Course, enrollments []models.Enrollment) (*EnrollmentCache, error) {
	cache := &EnrollmentCache{
		CourseCapacity:        make(map[uint]int),
		ConflictGraph:         make(map[uint]map[uint]bool),
		StudentCourses:        make(map[uint]map[uint]struct{}),
		StudentWaitingCourses: make(map[uint]map[uint]struct{}),
		EnrolledCount:         make(map[uint]*atomic.Int32),
		WaitingCount:          make(map[uint]*atomic.Int32),
	}
	cache.loadInitStudents(students)
	cache.loadInitCourses(courses)
	cache.loadEnrollments(enrollments)
	if err := cache.buildConflictGraph(courses); err != nil {
		return nil, err
	}
	return cache, nil
}

func (cache *EnrollmentCache) loadInitStudents(students []models.Student) {
	for _, s := range students {
		cache.StudentCourses[s.ID] = make(map[uint]struct{})
		cache.StudentWaitingCourses[s.ID] = make(map[uint]struct{})
	}
}

func (cache *EnrollmentCache) loadInitCourses(courses []models.Course) {
	for _, c := range courses {
		cache.CourseCapacity[c.ID] = c.Capacity
		cache.EnrolledCount[c.ID] = &atomic.Int32{}
		cache.WaitingCount[c.ID] = &atomic.Int32{}
	}
}

// loadEnrollments loads existing enrollments into cache
// Must be called after LoadInitStudents and LoadInitCourses
func (cache *EnrollmentCache) loadEnrollments(enrollments []models.Enrollment) {
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

func (cache *EnrollmentCache) buildConflictGraph(courses []models.Course) error {
	cache.ConflictGraph = make(map[uint]map[uint]bool)
	for i, course1 := range courses {
		cache.ConflictGraph[course1.ID] = make(map[uint]bool)
		for j, course2 := range courses {
			if i != j {
				conflict, err := hasCourseScheduleConflict(course1.Schedules, course2.Schedules)
				if err != nil {
					return fmt.Errorf("failed to check conflict between course %d and %d: %w", course1.ID, course2.ID, err)
				}
				if conflict {
					cache.ConflictGraph[course1.ID][course2.ID] = true
				}
			}
		}
	}
	return nil
}

// CourseExists checks if a course exists in cache
func (cache *EnrollmentCache) CourseExists(courseID uint) bool {
	_, exists := cache.CourseCapacity[courseID]
	return exists
}

// StudentExists checks if a student exists in cache
func (cache *EnrollmentCache) StudentExists(studentID uint) bool {
	return cache.StudentCourses[studentID] != nil
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

// IsStudentEnrolled checks if a student is already enrolled in a course
// Assumes student existence is already validated
func (cache *EnrollmentCache) IsStudentEnrolled(studentID, courseID uint) bool {
	_, exists := cache.StudentCourses[studentID][courseID]
	return exists
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

func (cache *EnrollmentCache) GetPosIfNotFull(courseID uint) (int, error) {
	capacity := cache.CourseCapacity[courseID]
	enrolledCount := int(cache.EnrolledCount[courseID].Load())
	if enrolledCount >= capacity {
		return 0, errors.New("")
	}
	return enrolledCount, nil
}

// IsWaitlistFull checks if a course's waitlist has reached capacity
// Assumes course existence is already validated
func (cache *EnrollmentCache) IsWaitlistFull(courseID uint) bool {
	capacity := cache.CourseCapacity[courseID]
	waitingCount := int(cache.WaitingCount[courseID].Load())
	return waitingCount >= capacity
}

// EnrollStudent enrolls a student in a course
// Assumes student and course existence is already validated
func (cache *EnrollmentCache) EnrollStudent(studentID, courseID uint) {
	cache.EnrolledCount[courseID].Add(1)
	cache.StudentCourses[studentID][courseID] = struct{}{}
}

// AddToWaitlist adds a student to a course's waitlist and returns their position
// Assumes student and course existence is already validated
func (cache *EnrollmentCache) AddToWaitlist(studentID, courseID uint) int {
	newCount := cache.WaitingCount[courseID].Add(1)
	cache.StudentWaitingCourses[studentID][courseID] = struct{}{}
	return int(newCount)
}
