package cache

import (
	"course-reg/internal/app/models"
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
	return &EnrollmentCache{
		Courses:          make(map[uint]*models.Course),
		ConflictGraph:    make(map[uint]map[uint]bool),
		EnrolledStudents: make(map[uint][]uint),
		WaitingStudents:  make(map[uint][]uint),
		StudentCourses:   make(map[uint]map[uint]bool),
	}
}

// LoadCourses loads courses into cache
func (c *EnrollmentCache) LoadCourses(courses []models.Course) {
	for i := range courses {
		c.Courses[courses[i].ID] = &courses[i]
		if _, exists := c.EnrolledStudents[courses[i].ID]; !exists {
			c.EnrolledStudents[courses[i].ID] = []uint{}
			c.WaitingStudents[courses[i].ID] = []uint{}
		}
	}
}

// LoadEnrollments loads existing enrollments into cache
func (c *EnrollmentCache) LoadEnrollments(enrollments []models.Enrollment) {
	for _, e := range enrollments {
		if e.IsWaitlist {
			c.WaitingStudents[e.CourseID] = append(c.WaitingStudents[e.CourseID], e.StudentID)
		} else {
			c.EnrolledStudents[e.CourseID] = append(c.EnrolledStudents[e.CourseID], e.StudentID)
		}

		if c.StudentCourses[e.StudentID] == nil {
			c.StudentCourses[e.StudentID] = make(map[uint]bool)
		}
		c.StudentCourses[e.StudentID][e.CourseID] = true
	}
}

// BuildConflictGraph builds conflict graph from courses
func (c *EnrollmentCache) BuildConflictGraph() {
	c.ConflictGraph = make(map[uint]map[uint]bool)

	for id1, course1 := range c.Courses {
		c.ConflictGraph[id1] = make(map[uint]bool)
		for id2, course2 := range c.Courses {
			if id1 != id2 && DoSchedulesConflict(course1.Schedules, course2.Schedules) {
				c.ConflictGraph[id1][id2] = true
			}
		}
	}
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
		if id != course.ID && DoSchedulesConflict(course.Schedules, existingCourse.Schedules) {
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
