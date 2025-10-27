package cache

import (
	"course-reg/internal/app/models"
	"course-reg/internal/pkg/constant"
	"course-reg/internal/pkg/util"
	"fmt"
	"sort"
)

// EnrollmentCache is a simple in-memory data structure
// NOTE: NOT thread-safe. Should only be accessed by a single goroutine (worker)
type EnrollmentCache struct {

	// Course data
	Courses       map[uint]*models.Course // courseID -> Course
	ConflictGraph map[uint]map[uint]bool  // courseID -> conflicting courseIDs

	// Enrollment data
	EnrolledStudents map[uint][]uint        // courseID -> []studentID
	WaitingStudents  map[uint][]uint        // courseID -> []studentID (waitlist)
	StudentCourses   map[uint]map[uint]bool // studentID -> set of courseIDs
}

func NewEnrollmentCache() *EnrollmentCache {
	// todo: init cache

	return &EnrollmentCache{
		Courses:          make(map[uint]*models.Course),
		ConflictGraph:    make(map[uint]map[uint]bool),
		EnrolledStudents: make(map[uint][]uint),
		WaitingStudents:  make(map[uint][]uint),
		StudentCourses:   make(map[uint]map[uint]bool),
	}
}

// LoadEnrollments loads existing enrollments into cache
// func (c *EnrollmentCache) LoadEnrollments(enrollments []models.Enrollment) {
// 	for _, e := range enrollments {
// 		if e.IsWaitlist {
// 			c.WaitingStudents[e.CourseID] = append(c.WaitingStudents[e.CourseID], e.StudentID)
// 		} else {
// 			c.EnrolledStudents[e.CourseID] = append(c.EnrolledStudents[e.CourseID], e.StudentID)
// 		}

// 		if c.StudentCourses[e.StudentID] == nil {
// 			c.StudentCourses[e.StudentID] = make(map[uint]bool)
// 		}
// 		c.StudentCourses[e.StudentID][e.CourseID] = true
// 	}
// }

func (c *EnrollmentCache) LoadInitStudents(students []models.Student) {
	for _, s := range students {
		c.StudentCourses[s.ID] = make(map[uint]bool)
	}
}

func (c *EnrollmentCache) LoadInitCourses(courses []models.Course) {
	for i := range courses {
		c.Courses[courses[i].ID] = &courses[i]
		if _, exists := c.EnrolledStudents[courses[i].ID]; !exists {
			c.EnrolledStudents[courses[i].ID] = []uint{}
			c.WaitingStudents[courses[i].ID] = []uint{}
		}
	}
}

func (c *EnrollmentCache) BuildConflictGraph() {
	c.ConflictGraph = make(map[uint]map[uint]bool)
	for id1, course1 := range c.Courses {
		c.ConflictGraph[id1] = make(map[uint]bool)
		for id2, course2 := range c.Courses {
			if id1 != id2 && util.SchedulesConflict(course1.Schedules, course2.Schedules) {
				c.ConflictGraph[id1][id2] = true
			}
		}
	}
}

// AddCourse adds a single course and updates conflict graph incrementally
func (c *EnrollmentCache) AddCourse(course models.Course) {
	// Add course to cache
	c.Courses[course.ID] = &course
	c.EnrolledStudents[course.ID] = []uint{}
	c.WaitingStudents[course.ID] = []uint{}

	// Initialize conflict graph entry
	c.ConflictGraph[course.ID] = make(map[uint]bool)

	// Check conflicts with all existing courses
	for id, existingCourse := range c.Courses {
		if id != course.ID && util.SchedulesConflict(course.Schedules, existingCourse.Schedules) {
			c.ConflictGraph[course.ID][id] = true
			c.ConflictGraph[id][course.ID] = true
		}
	}
}

// RemoveCourse removes a course and updates conflict graph
func (c *EnrollmentCache) RemoveCourse(courseID uint) {
	// Remove from cache
	delete(c.Courses, courseID)
	delete(c.EnrolledStudents, courseID)
	delete(c.WaitingStudents, courseID)

	// Remove from conflict graph
	delete(c.ConflictGraph, courseID)
	for id := range c.ConflictGraph {
		delete(c.ConflictGraph[id], courseID)
	}
}

// ClearAllCourses clears all courses from cache
func (c *EnrollmentCache) ClearAllCourses() {
	c.Courses = make(map[uint]*models.Course)
	c.EnrolledStudents = make(map[uint][]uint)
	c.WaitingStudents = make(map[uint][]uint)
	c.ConflictGraph = make(map[uint]map[uint]bool)
}

func (c *EnrollmentCache) ClearAllStudents() {
	c.StudentCourses = make(map[uint]map[uint]bool)
}

func (c *EnrollmentCache) GetAllCourseStatus() map[uint]constant.CourseStatus {
	status := make(map[uint]constant.CourseStatus)
	for courseID, course := range c.Courses {
		if len(c.EnrolledStudents[courseID]) < course.Capacity {
			status[courseID] = constant.CourseAvailable
		} else if len(c.WaitingStudents[courseID]) < course.Capacity {
			status[courseID] = constant.CourseWaitlist
		} else {
			status[courseID] = constant.CourseFull
		}
	}
	return status
}

// ========== Business Logic Methods for Testing ==========

// GetCourse retrieves a course by ID
func (c *EnrollmentCache) GetCourse(courseID uint) (*models.Course, bool) {
	course, exists := c.Courses[courseID]
	return course, exists
}

// IsStudentEnrolled checks if a student is already enrolled in a course
func (c *EnrollmentCache) IsStudentEnrolled(studentID, courseID uint) bool {
	return c.StudentCourses[studentID] != nil && c.StudentCourses[studentID][courseID]
}

// HasTimeConflict checks if enrolling in a course would create a time conflict
func (c *EnrollmentCache) HasTimeConflict(studentID, courseID uint) bool {
	if c.StudentCourses[studentID] == nil {
		return false
	}
	for enrolledCourse := range c.StudentCourses[studentID] {
		if c.ConflictGraph[courseID][enrolledCourse] {
			return true
		}
	}
	return false
}

// IsFull checks if a course has reached its enrollment capacity
func (c *EnrollmentCache) IsFull(courseID uint) bool {
	course, exists := c.Courses[courseID]
	if !exists {
		return true
	}
	return len(c.EnrolledStudents[courseID]) >= course.Capacity
}

// IsWaitlistFull checks if a course's waitlist has reached capacity
func (c *EnrollmentCache) IsWaitlistFull(courseID uint) bool {
	course, exists := c.Courses[courseID]
	if !exists {
		return true
	}
	return len(c.WaitingStudents[courseID]) >= course.Capacity
}

// EnrollStudent enrolls a student in a course
func (c *EnrollmentCache) EnrollStudent(studentID, courseID uint) {
	c.EnrolledStudents[courseID] = append(c.EnrolledStudents[courseID], studentID)
	if c.StudentCourses[studentID] == nil {
		c.StudentCourses[studentID] = make(map[uint]bool)
	}
	c.StudentCourses[studentID][courseID] = true
}

// AddToWaitlist adds a student to a course's waitlist and returns their position
func (c *EnrollmentCache) AddToWaitlist(studentID, courseID uint) int {
	c.WaitingStudents[courseID] = append(c.WaitingStudents[courseID], studentID)
	return len(c.WaitingStudents[courseID])
}

func (c *EnrollmentCache) DebugPrint() {
	fmt.Println("========== EnrollmentCache Debug ==========")

	// ----- Courses -----
	fmt.Println("📘 Courses:")
	for id, course := range c.Courses {
		fmt.Printf("  - [%d] %s (Instructor: %s)\n", id, course.Name, course.Instructor)
	}
	if len(c.Courses) == 0 {
		fmt.Println("  (no courses loaded)")
	}

	// ----- ConflictGraph -----
	fmt.Println("\n⚔️  ConflictGraph:")
	for courseID, conflicts := range c.ConflictGraph {
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
	if len(c.ConflictGraph) == 0 {
		fmt.Println("  (no conflicts registered)")
	}

	// ----- EnrolledStudents -----
	fmt.Println("\n👩‍🎓 EnrolledStudents:")
	for courseID, students := range c.EnrolledStudents {
		fmt.Printf("  - Course %d: %v\n", courseID, students)
	}
	if len(c.EnrolledStudents) == 0 {
		fmt.Println("  (no enrolled students)")
	}

	// ----- WaitingStudents -----
	fmt.Println("\n🕒 WaitingStudents:")
	for courseID, students := range c.WaitingStudents {
		fmt.Printf("  - Course %d: %v\n", courseID, students)
	}
	if len(c.WaitingStudents) == 0 {
		fmt.Println("  (no waiting students)")
	}

	// ----- StudentCourses -----
	fmt.Println("\n📚 StudentCourses:")
	for studentID, courseSet := range c.StudentCourses {
		courseIDs := make([]uint, 0, len(courseSet))
		for cid := range courseSet {
			courseIDs = append(courseIDs, cid)
		}
		sort.Slice(courseIDs, func(i, j int) bool { return courseIDs[i] < courseIDs[j] })
		fmt.Printf("  - Student %d: %v\n", studentID, courseIDs)
	}
	if len(c.StudentCourses) == 0 {
		fmt.Println("  (no student-course mappings)")
	}

	fmt.Println("===========================================")
	fmt.Println()
}
