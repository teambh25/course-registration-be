package cache

import (
	"course-reg/internal/app/models"
	"course-reg/internal/pkg/constant"
	util "course-reg/internal/pkg/utils"
	"errors"
	"fmt"
	"sort"
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
	EnrolledSeq           map[uint]*atomic.Int32     // courseID -> count of enrolled students (atomic)
	WaitingSeq            map[uint]*atomic.Int32     // courseID -> count of waiting students (atomic)
}

func NewEnrollmentCache() *EnrollmentCache {
	return &EnrollmentCache{
		CourseCapacity:        make(map[uint]int),
		ConflictGraph:         make(map[uint]map[uint]bool),
		StudentCourses:        make(map[uint]map[uint]struct{}),
		StudentWaitingCourses: make(map[uint]map[uint]struct{}),
		EnrolledSeq:           make(map[uint]*atomic.Int32),
		WaitingSeq:            make(map[uint]*atomic.Int32),
	}
}

func NewEnrollmentCacheWithData(students []models.Student, courses []models.Course, enrollments []models.Enrollment) *EnrollmentCache {
	cache := NewEnrollmentCache()
	cache.loadInitStudents(students)
	cache.loadInitCourses(courses)
	cache.loadEnrollments(enrollments)
	cache.buildConflictGraph(courses)
	return cache
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
		cache.EnrolledSeq[c.ID] = &atomic.Int32{}
		cache.WaitingSeq[c.ID] = &atomic.Int32{}
	}
}

// loadEnrollments loads existing enrollments into cache
// Must be called after LoadInitStudents and LoadInitCourses
func (cache *EnrollmentCache) loadEnrollments(enrollments []models.Enrollment) {
	for _, e := range enrollments {
		if e.IsWaitlist {
			// Update to max(current, position + 1)
			currentMax := cache.WaitingSeq[e.CourseID].Load()
			if int32(e.Position+1) > currentMax {
				cache.WaitingSeq[e.CourseID].Store(int32(e.Position + 1))
			}
			cache.StudentWaitingCourses[e.StudentID][e.CourseID] = struct{}{}
		} else {
			// Update to max(current, position + 1)
			currentMax := cache.EnrolledSeq[e.CourseID].Load()
			if int32(e.Position+1) > currentMax {
				cache.EnrolledSeq[e.CourseID].Store(int32(e.Position + 1))
			}
			cache.StudentCourses[e.StudentID][e.CourseID] = struct{}{}
		}
	}
}

func (cache *EnrollmentCache) buildConflictGraph(courses []models.Course) {
	cache.ConflictGraph = make(map[uint]map[uint]bool)
	for i, course1 := range courses {
		cache.ConflictGraph[course1.ID] = make(map[uint]bool)
		for j, course2 := range courses {
			if i != j && util.SchedulesConflict(course1.Schedules, course2.Schedules) {
				cache.ConflictGraph[course1.ID][course2.ID] = true
			}
		}
	}
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

func (cache *EnrollmentCache) GetAllCourseStatus() map[uint]constant.CourseStatus {
	status := make(map[uint]constant.CourseStatus)
	for courseID, capacity := range cache.CourseCapacity {
		enrolledCount := int(cache.EnrolledSeq[courseID].Load())
		waitingCount := int(cache.WaitingSeq[courseID].Load())

		if enrolledCount < capacity {
			status[courseID] = constant.CourseAvailable
		} else if waitingCount < capacity {
			status[courseID] = constant.CourseWaitlist
		} else {
			status[courseID] = constant.CourseFull
		}
	}
	return status
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
	enrolledCount := int(cache.EnrolledSeq[courseID].Load())
	if enrolledCount >= capacity {
		return 0, errors.New("")
	}
	return enrolledCount, nil
}

// IsWaitlistFull checks if a course's waitlist has reached capacity
// Assumes course existence is already validated
func (cache *EnrollmentCache) IsWaitlistFull(courseID uint) bool {
	capacity := cache.CourseCapacity[courseID]
	waitingCount := int(cache.WaitingSeq[courseID].Load())
	return waitingCount >= capacity
}

// EnrollStudent enrolls a student in a course
// Assumes student and course existence is already validated
func (cache *EnrollmentCache) EnrollStudent(studentID, courseID uint) {
	cache.EnrolledSeq[courseID].Add(1)
	cache.StudentCourses[studentID][courseID] = struct{}{}
}

// AddToWaitlist adds a student to a course's waitlist and returns their position
// Assumes student and course existence is already validated
func (cache *EnrollmentCache) AddToWaitlist(studentID, courseID uint) int {
	newCount := cache.WaitingSeq[courseID].Add(1)
	cache.StudentWaitingCourses[studentID][courseID] = struct{}{}
	return int(newCount)
}

func (cache *EnrollmentCache) DebugPrint() {
	fmt.Println("========== EnrollmentCache Debug ==========")

	// ----- Course Capacity -----
	fmt.Println("📘 Course Capacity:")
	for courseID, capacity := range cache.CourseCapacity {
		enrolledCount := int(cache.EnrolledSeq[courseID].Load())
		waitingCount := int(cache.WaitingSeq[courseID].Load())
		fmt.Printf("  - Course %d: Capacity=%d, Enrolled=%d, Waiting=%d\n",
			courseID, capacity, enrolledCount, waitingCount)
	}
	if len(cache.CourseCapacity) == 0 {
		fmt.Println("  (no courses loaded)")
	}

	// ----- ConflictGraph -----
	fmt.Println("\n⚔️  ConflictGraph:")
	for courseID, conflicts := range cache.ConflictGraph {
		if len(conflicts) == 0 {
			continue
		}
		conflictIDs := make([]uint, 0, len(conflicts))
		for cid := range conflicts {
			conflictIDs = append(conflictIDs, cid)
		}
		sort.Slice(conflictIDs, func(i, j int) bool { return conflictIDs[i] < conflictIDs[j] })
		fmt.Printf("  - Course %d conflicts with %v\n", courseID, conflictIDs)
	}
	if len(cache.ConflictGraph) == 0 {
		fmt.Println("  (no conflicts registered)")
	}

	// ----- StudentCourses -----
	fmt.Println("\n📚 StudentCourses (Enrolled):")
	for studentID, courseSet := range cache.StudentCourses {
		courseIDs := make([]uint, 0, len(courseSet))
		for cid := range courseSet {
			courseIDs = append(courseIDs, cid)
		}
		sort.Slice(courseIDs, func(i, j int) bool { return courseIDs[i] < courseIDs[j] })
		fmt.Printf("  - Student %d: %v\n", studentID, courseIDs)
	}
	if len(cache.StudentCourses) == 0 {
		fmt.Println("  (no student-course mappings)")
	}

	// ----- StudentWaitingCourses -----
	fmt.Println("\n🕒 StudentWaitingCourses:")
	for studentID, courseSet := range cache.StudentWaitingCourses {
		if len(courseSet) == 0 {
			continue
		}
		courseIDs := make([]uint, 0, len(courseSet))
		for cid := range courseSet {
			courseIDs = append(courseIDs, cid)
		}
		sort.Slice(courseIDs, func(i, j int) bool { return courseIDs[i] < courseIDs[j] })
		fmt.Printf("  - Student %d: %v\n", studentID, courseIDs)
	}

	fmt.Println("===========================================")
	fmt.Println()
}
