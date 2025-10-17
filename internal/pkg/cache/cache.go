package cache

import (
	"course-reg/internal/app/models"
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
	// todo : init cache

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

	c.DebugPrint()
}

func (c *EnrollmentCache) LoadInitCourses(courses []models.Course) {
	for i := range courses {
		c.Courses[courses[i].ID] = &courses[i]
		if _, exists := c.EnrolledStudents[courses[i].ID]; !exists {
			c.EnrolledStudents[courses[i].ID] = []uint{}
			c.WaitingStudents[courses[i].ID] = []uint{}
		}
	}
	c.BuildConflictGraph()

	c.DebugPrint()
}

func (c *EnrollmentCache) BuildConflictGraph() {
	c.ConflictGraph = make(map[uint]map[uint]bool)

	for id1, course1 := range c.Courses {
		c.ConflictGraph[id1] = make(map[uint]bool)
		for id2, course2 := range c.Courses {
			if id1 != id2 && SchedulesConflict(course1.Schedules, course2.Schedules) {
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
		if id != course.ID && SchedulesConflict(course.Schedules, existingCourse.Schedules) {
			c.ConflictGraph[course.ID][id] = true
			c.ConflictGraph[id][course.ID] = true
		}
	}

	c.DebugPrint()
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

	c.DebugPrint()
}

// ClearAllCourses clears all courses from cache
func (c *EnrollmentCache) ClearAllCourses() {
	c.Courses = make(map[uint]*models.Course)
	c.EnrolledStudents = make(map[uint][]uint)
	c.WaitingStudents = make(map[uint][]uint)
	c.ConflictGraph = make(map[uint]map[uint]bool)

	c.DebugPrint()
}

func (c *EnrollmentCache) ClearAllStudents() {
	c.StudentCourses = make(map[uint]map[uint]bool)

	c.DebugPrint()
}

// GetAllRemainingSeats returns map of courseID -> remaining seats
func (c *EnrollmentCache) GetAllRemainingSeats() map[uint]int {
	result := make(map[uint]int)
	for courseID, course := range c.Courses {
		enrolled := len(c.EnrolledStudents[courseID])
		result[courseID] = course.Capacity - enrolled
	}
	return result
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

	fmt.Println("===========================================\n")
}
