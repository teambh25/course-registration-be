package cache

import (
	"course-reg/internal/app/models"
	"course-reg/internal/pkg/constant"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Helper Functions
func createTestCourses() []models.Course {
	return []models.Course{
		{ID: 1, Name: "데이터구조", Instructor: "김교수", Schedules: "월 09:00~11:00", Capacity: 30},
		{ID: 2, Name: "알고리즘", Instructor: "이교수", Schedules: "월 10:00~12:00", Capacity: 25},
		{ID: 3, Name: "운영체제", Instructor: "박교수", Schedules: "수 14:00~16:00", Capacity: 20},
	}
}

func createTestStudents() []models.Student {
	return []models.Student{
		{ID: 1, Name: "학생1", PhoneNumber: "010-1111-1111", BirthDate: "2000-01-01"},
		{ID: 2, Name: "학생2", PhoneNumber: "010-2222-2222", BirthDate: "2000-02-02"},
		{ID: 3, Name: "학생3", PhoneNumber: "010-3333-3333", BirthDate: "2000-03-03"},
	}
}

// Initialization Tests
func TestNewEnrollmentCache(t *testing.T) {
	cache := NewEnrollmentCache()

	assert.NotNil(t, cache, "NewEnrollmentCache returned nil")
	assert.NotNil(t, cache.Courses, "Courses map not initialized")
	assert.NotNil(t, cache.ConflictGraph, "ConflictGraph map not initialized")
	assert.NotNil(t, cache.EnrolledStudents, "EnrolledStudents map not initialized")
	assert.NotNil(t, cache.WaitingStudents, "WaitingStudents map not initialized")
	assert.NotNil(t, cache.StudentCourses, "StudentCourses map not initialized")
}

func TestLoadInitStudents(t *testing.T) {
	cache := NewEnrollmentCache()
	students := createTestStudents()

	cache.LoadInitStudents(students)

	for _, student := range students {
		assert.NotNil(t, cache.StudentCourses[student.ID], "Student %d not initialized in StudentCourses", student.ID)
		assert.Equal(t, 0, len(cache.StudentCourses[student.ID]), "Student %d should have no courses initially", student.ID)
	}
}

func TestLoadInitCourses(t *testing.T) {
	cache := NewEnrollmentCache()
	courses := createTestCourses()

	cache.LoadInitCourses(courses)

	// Check all courses are loaded
	assert.Equal(t, len(courses), len(cache.Courses), "Expected %d courses", len(courses))

	// Check each course
	for _, course := range courses {
		loaded, exists := cache.Courses[course.ID]
		assert.True(t, exists, "Course %d not loaded", course.ID)
		if !exists {
			continue
		}
		assert.Equal(t, course.Name, loaded.Name, "Course %d name mismatch", course.ID)

		// Check enrollment lists initialized
		assert.NotNil(t, cache.EnrolledStudents[course.ID], "EnrolledStudents not initialized for course %d", course.ID)
		assert.NotNil(t, cache.WaitingStudents[course.ID], "WaitingStudents not initialized for course %d", course.ID)
	}
}

// ========== Conflict Graph Tests ==========

func TestBuildConflictGraph(t *testing.T) {
	cache := NewEnrollmentCache()
	courses := []models.Course{
		{ID: 1, Name: "Course1", Schedules: "월 09:00~11:00, 목 09:00~11:00", Capacity: 30}, // Conflicts with Course 2 & 3
		{ID: 2, Name: "Course2", Schedules: "월 10:00~12:00", Capacity: 30},                // Conflicts with Course 1
		{ID: 3, Name: "Course3", Schedules: "수 14:00~16:00, 목 10:30~12:00", Capacity: 30}, // Conflicts with Course 1
		{ID: 4, Name: "Course4", Schedules: "금 14:00~16:00", Capacity: 30},                // No conflict
	}

	cache.LoadInitCourses(courses)
	cache.BuildConflictGraph()

	// Course 1
	assert.False(t, cache.ConflictGraph[1][1], "Course 1 should not conflict with itself")
	assert.True(t, cache.ConflictGraph[1][2], "Course 1 and 2 should conflict")
	assert.True(t, cache.ConflictGraph[1][3], "Course 1 and 3 should conflict")
	assert.False(t, cache.ConflictGraph[1][4], "Course 1 and 4 should not conflict")

	// Course 2
	assert.True(t, cache.ConflictGraph[2][1], "Course 2 and 1 should conflict")
	assert.False(t, cache.ConflictGraph[2][2], "Course 2 should not conflict with itself")
	assert.False(t, cache.ConflictGraph[2][3], "Course 2 and 3 should not conflict")
	assert.False(t, cache.ConflictGraph[2][4], "Course 2 and 4 should not conflict")

	// Course 3
	assert.True(t, cache.ConflictGraph[3][1], "Course 3 and 1 should conflict")
	assert.False(t, cache.ConflictGraph[3][2], "Course 3 and 2 should not conflict")
	assert.False(t, cache.ConflictGraph[3][3], "Course 3 should not conflict with itself")
	assert.False(t, cache.ConflictGraph[3][4], "Course 3 and 4 should not conflict")

	// Course 4
	assert.False(t, cache.ConflictGraph[4][1], "Course 4 and 1 should not conflict")
	assert.False(t, cache.ConflictGraph[4][2], "Course 4 and 2 should not conflict")
	assert.False(t, cache.ConflictGraph[4][3], "Course 4 and 3 should not conflict")
	assert.False(t, cache.ConflictGraph[4][4], "Course 4 should not conflict with itself")
}

// ========== Course Management Tests ==========

func TestAddCourse(t *testing.T) {
	cache := NewEnrollmentCache()
	initialCourses := []models.Course{
		{ID: 1, Name: "Course1", Schedules: "월 09:00~11:00", Capacity: 30},
	}
	cache.LoadInitCourses(initialCourses)
	cache.BuildConflictGraph()

	newCourse := models.Course{ID: 2, Name: "Course2", Schedules: "월 10:00~12:00", Capacity: 25}
	cache.AddCourse(newCourse)

	// Check course is added
	course, exists := cache.Courses[2]
	assert.True(t, exists, "Course 2 not added")
	assert.Equal(t, "Course2", course.Name, "Course name mismatch")

	// Check enrollment lists initialized
	assert.NotNil(t, cache.EnrolledStudents[2], "EnrolledStudents not initialized for new course")
	assert.NotNil(t, cache.WaitingStudents[2], "WaitingStudents not initialized for new course")

	// Check conflict graph updated
	assert.True(t, cache.ConflictGraph[1][2], "Conflict graph not updated: Course 1 should conflict with Course 2")
	assert.True(t, cache.ConflictGraph[2][1], "Conflict graph not updated: Course 2 should conflict with Course 1")
	assert.False(t, cache.ConflictGraph[2][2], "Conflict graph not updated: Course 2 should not conflict with itself")
}

func TestRemoveCourse(t *testing.T) {
	cache := NewEnrollmentCache()
	courses := []models.Course{
		{ID: 1, Name: "Course1", Schedules: "월 09:00~11:00", Capacity: 30},
		{ID: 2, Name: "Course2", Schedules: "수 14:00~16:00", Capacity: 25},
	}
	cache.LoadInitCourses(courses)

	cache.RemoveCourse(1)

	// Check course is removed
	_, exists := cache.Courses[1]
	assert.False(t, exists, "Course 1 should be removed")

	// Check enrollment lists removed
	_, exists = cache.EnrolledStudents[1]
	assert.False(t, exists, "EnrolledStudents for course 1 should be removed")
	_, exists = cache.WaitingStudents[1]
	assert.False(t, exists, "WaitingStudents for course 1 should be removed")

	// Check conflict graph cleaned up
	_, exists = cache.ConflictGraph[1]
	assert.False(t, exists, "ConflictGraph entry for course 1 should be removed")
	assert.False(t, cache.ConflictGraph[2][1], "ConflictGraph references to course 1 should be removed")
}

func TestClearAllCourses(t *testing.T) {
	cache := NewEnrollmentCache()
	cache.LoadInitCourses(createTestCourses())

	cache.ClearAllCourses()

	assert.Equal(t, 0, len(cache.Courses), "Expected 0 courses after clear")
	assert.Equal(t, 0, len(cache.EnrolledStudents), "Expected 0 enrolled student entries after clear")
	assert.Equal(t, 0, len(cache.WaitingStudents), "Expected 0 waiting student entries after clear")
	assert.Equal(t, 0, len(cache.ConflictGraph), "Expected 0 conflict graph entries after clear")
}

func TestClearAllStudents(t *testing.T) {
	cache := NewEnrollmentCache()
	cache.LoadInitStudents(createTestStudents())

	cache.ClearAllStudents()

	assert.Equal(t, 0, len(cache.StudentCourses), "Expected 0 student course entries after clear")
}

// ========== Business Logic Method Tests ==========

func TestGetCourse(t *testing.T) {
	cache := NewEnrollmentCache()
	courses := createTestCourses()
	cache.LoadInitCourses(courses)

	// Test existing course
	course, exists := cache.GetCourse(1)
	assert.True(t, exists, "Course 1 should exist")
	assert.Equal(t, "데이터구조", course.Name, "Course name mismatch")

	// Test non-existent course
	_, exists = cache.GetCourse(999)
	assert.False(t, exists, "Course 999 should not exist")
}

func TestIsStudentEnrolled(t *testing.T) {
	cache := NewEnrollmentCache()
	cache.LoadInitStudents(createTestStudents())
	cache.LoadInitCourses(createTestCourses())

	// Initially not enrolled
	assert.False(t, cache.IsStudentEnrolled(1, 1), "Student 1 should not be enrolled in course 1 initially")

	// Enroll student
	cache.EnrollStudent(1, 1)

	// Now enrolled
	assert.True(t, cache.IsStudentEnrolled(1, 1), "Student 1 should be enrolled in course 1 after enrollment")

	// Still not enrolled in other course
	assert.False(t, cache.IsStudentEnrolled(1, 2), "Student 1 should not be enrolled in course 2")
}

func TestHasTimeConflict(t *testing.T) {
	cache := NewEnrollmentCache()
	cache.LoadInitStudents(createTestStudents())
	courses := []models.Course{
		{ID: 1, Name: "Course1", Schedules: "월 09:00~11:00", Capacity: 30},
		{ID: 2, Name: "Course2", Schedules: "월 10:00~12:00", Capacity: 30}, // Conflicts with Course1
		{ID: 3, Name: "Course3", Schedules: "수 14:00~16:00", Capacity: 30}, // No conflict
	}
	cache.LoadInitCourses(courses)
	cache.BuildConflictGraph()

	// No conflict initially
	assert.False(t, cache.HasTimeConflict(1, 1), "No conflict expected for student 1 enrolling in course 1 (first course)")

	// Enroll in Course1
	cache.EnrollStudent(1, 1)

	// Should have conflict with Course2
	assert.True(t, cache.HasTimeConflict(1, 2), "Expected time conflict between Course 1 and Course 2")

	// Should not have conflict with Course3
	assert.False(t, cache.HasTimeConflict(1, 3), "No conflict expected between Course 1 and Course 3")
}

func TestIsFull(t *testing.T) {
	cache := NewEnrollmentCache()
	cache.LoadInitStudents(createTestStudents())
	courses := []models.Course{
		{ID: 1, Name: "Course1", Schedules: "월 09:00~11:00", Capacity: 2},
	}
	cache.LoadInitCourses(courses)

	// Initially not full
	assert.False(t, cache.IsFull(1), "Course 1 should not be full initially")

	// Enroll 2 students (fill to capacity)
	cache.EnrollStudent(1, 1)
	cache.EnrollStudent(2, 1)

	// Now full
	assert.True(t, cache.IsFull(1), "Course 1 should be full after enrolling 2 students")

	// Non-existent course should return true (full)
	assert.True(t, cache.IsFull(999), "Non-existent course should be considered full")
}

func TestIsWaitlistFull(t *testing.T) {
	cache := NewEnrollmentCache()
	cache.LoadInitStudents(createTestStudents())
	courses := []models.Course{
		{ID: 1, Name: "Course1", Schedules: "월 09:00~11:00", Capacity: 2},
	}
	cache.LoadInitCourses(courses)

	// Initially not full
	assert.False(t, cache.IsWaitlistFull(1), "Waitlist should not be full initially")

	// Add 2 students to waitlist (fill to capacity)
	cache.AddToWaitlist(1, 1)
	cache.AddToWaitlist(2, 1)

	// Now full
	assert.True(t, cache.IsWaitlistFull(1), "Waitlist should be full after adding 2 students")

	// Non-existent course should return true (full)
	assert.True(t, cache.IsWaitlistFull(999), "Non-existent course waitlist should be considered full")
}

func TestEnrollStudent(t *testing.T) {
	cache := NewEnrollmentCache()
	cache.LoadInitStudents(createTestStudents())
	cache.LoadInitCourses(createTestCourses())

	cache.EnrollStudent(1, 1)

	// Check student is in enrolled list
	found := false
	for _, sid := range cache.EnrolledStudents[1] {
		if sid == 1 {
			found = true
			break
		}
	}
	assert.True(t, found, "Student 1 should be in enrolled students list for course 1")

	// Check student courses mapping
	assert.True(t, cache.StudentCourses[1][1], "Course 1 should be in student 1's course set")
}

func TestAddToWaitlist(t *testing.T) {
	cache := NewEnrollmentCache()
	cache.LoadInitStudents(createTestStudents())
	cache.LoadInitCourses(createTestCourses())

	// Add first student to waitlist
	position1 := cache.AddToWaitlist(1, 1)
	assert.Equal(t, 1, position1, "Expected waitlist position 1")

	// Add second student to waitlist
	position2 := cache.AddToWaitlist(2, 1)
	assert.Equal(t, 2, position2, "Expected waitlist position 2")

	// Check students are in waitlist
	assert.Equal(t, 2, len(cache.WaitingStudents[1]), "Expected 2 students in waitlist")
}

// ========== Course Status Tests ==========

func TestGetAllCourseStatus(t *testing.T) {
	cache := NewEnrollmentCache()
	students := make([]models.Student, 5)
	for i := 0; i < 5; i++ {
		students[i] = models.Student{
			ID:          uint(i + 1),
			Name:        "Student",
			PhoneNumber: "010-0000-0000",
			BirthDate:   "2000-01-01",
		}
	}
	cache.LoadInitStudents(students)

	courses := []models.Course{
		{ID: 1, Name: "Course1", Schedules: "월 09:00~11:00", Capacity: 2},
		{ID: 2, Name: "Course2", Schedules: "수 14:00~16:00", Capacity: 2},
		{ID: 3, Name: "Course3", Schedules: "금 10:00~12:00", Capacity: 2},
	}
	cache.LoadInitCourses(courses)

	// Test AVAILABLE status
	statuses := cache.GetAllCourseStatus()
	assert.Equal(t, constant.CourseAvailable, statuses[1], "Course 1 should be AVAILABLE")

	// Fill course 1 (should become WAITLIST)
	cache.EnrollStudent(1, 1)
	cache.EnrollStudent(2, 1)
	statuses = cache.GetAllCourseStatus()
	assert.Equal(t, constant.CourseWaitlist, statuses[1], "Course 1 should be WAITLIST after filling")

	// Fill waitlist for course 1 (should become FULL)
	cache.AddToWaitlist(3, 1)
	cache.AddToWaitlist(4, 1)
	statuses = cache.GetAllCourseStatus()
	assert.Equal(t, constant.CourseFull, statuses[1], "Course 1 should be FULL after filling waitlist")
}

// ========== Edge Case Tests ==========

func TestEnrollStudent_NewStudent(t *testing.T) {
	cache := NewEnrollmentCache()
	cache.LoadInitCourses(createTestCourses())

	// Enroll a student that wasn't in LoadInitStudents
	cache.EnrollStudent(999, 1)

	// Should initialize StudentCourses for new student
	assert.NotNil(t, cache.StudentCourses[999], "StudentCourses should be initialized for new student")
	assert.True(t, cache.StudentCourses[999][1], "Course 1 should be in new student's course set")
}

func TestIsStudentEnrolled_NonExistentStudent(t *testing.T) {
	cache := NewEnrollmentCache()
	cache.LoadInitCourses(createTestCourses())

	// Should return false for non-existent student
	assert.False(t, cache.IsStudentEnrolled(999, 1), "Non-existent student should not be enrolled")
}

func TestHasTimeConflict_NonExistentStudent(t *testing.T) {
	cache := NewEnrollmentCache()
	cache.LoadInitCourses(createTestCourses())

	// Should return false for non-existent student
	assert.False(t, cache.HasTimeConflict(999, 1), "Non-existent student should not have time conflicts")
}
