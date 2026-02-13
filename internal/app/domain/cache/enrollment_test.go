package cache

import (
	"course-reg/internal/app/models"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Helper Functions

func createTestCourses() []models.Course {
	return []models.Course{
		{ID: 1, Name: "데이터구조", Instructor: "김교수", Schedules: "월 09:00~11:00", Capacity: 3},
		{ID: 2, Name: "알고리즘", Instructor: "이교수", Schedules: "월 10:00~12:00", Capacity: 2},
		{ID: 3, Name: "운영체제", Instructor: "박교수", Schedules: "수 14:00~16:00", Capacity: 2},
	}
}

func createTestStudents() []models.Student {
	return []models.Student{
		{ID: 1, Name: "학생1", PhoneNumber: "010-1111-1111", BirthDate: "1997-01-05"},
		{ID: 2, Name: "학생2", PhoneNumber: "010-2222-2222", BirthDate: "2000-06-15"},
		{ID: 3, Name: "학생3", PhoneNumber: "010-3333-3333", BirthDate: "2001-12-23"},
	}
}

func createTestEnrollments() []models.Enrollment {
	return []models.Enrollment{
		// Course 1 (0/3): empty
		// Course 2 (1/2): partially filled
		{StudentID: 1, CourseID: 2, Position: 0},

		// Course 3 (2/2): full
		{StudentID: 2, CourseID: 3, Position: 0},
		{StudentID: 3, CourseID: 3, Position: 1},
	}
}

func createTestCache(t *testing.T) *EnrollmentCache {
	t.Helper()
	c, err := NewEnrollmentCacheWithData(createTestStudents(), createTestCourses(), createTestEnrollments())
	assert.NoError(t, err)
	return c
}

// ========== Business Logic Tests ==========

func TestCourseExists(t *testing.T) {
	cache := createTestCache(t)
	existingID := uint(1)

	assert.True(t, cache.CourseExists(existingID))
	assert.False(t, cache.CourseExists(math.MaxUint32)) // non-existent course ID
}

func TestStudentExists(t *testing.T) {
	cache := createTestCache(t)
	existingID := uint(1)

	assert.True(t, cache.StudentExists(existingID))
	assert.False(t, cache.StudentExists(math.MaxUint32)) // non-existent student ID
}

func TestHasTimeConflict(t *testing.T) {
	cache := createTestCache(t) // Student 1 is enrolled in Course 2 (월 10:00~12:00)

	tests := []struct {
		name     string
		courseID uint
		expected bool
	}{
		{"overlapping schedule", 1, true},      // Course 1: 월 09:00~11:00 (conflicts with 2)
		{"non-overlapping schedule", 3, false}, // Course 3: 수 14:00~16:00 (no conflict with 2)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, cache.HasTimeConflict(1, tt.courseID))
		})
	}
}

func TestIsStudentEnrolled(t *testing.T) {
	cache := createTestCache(t) // Student 1 is enrolled in Course 2
	studentID := uint(1)
	enrolledCourseID := uint(2)
	notEnrolledCourseID := uint(1)

	assert.True(t, cache.IsStudentEnrolled(studentID, enrolledCourseID))
	assert.False(t, cache.IsStudentEnrolled(studentID, notEnrolledCourseID))
}

func TestGetAvailablePos(t *testing.T) {
	cache := createTestCache(t)

	tests := []struct {
		name        string
		courseID    uint
		expectedPos int
		expectedOk  bool
	}{
		{"empty course", 1, 0, true},     // Course 1 (0/3)
		{"partially filled", 2, 1, true}, // Course 2 (1/2)
		{"full course", 3, 0, false},     // Course 3 (2/2)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos, ok := cache.GetAvailablePos(tt.courseID)
			assert.Equal(t, tt.expectedOk, ok)
			if ok {
				assert.Equal(t, tt.expectedPos, pos)
			}
		})
	}
}

// 나중에 수정
// func TestGetAllCourseCountInfo(t *testing.T) {
// 	students := createTestStudents()
// 	courses := []models.Course{
// 		{ID: 1, Name: "Course1", Schedules: "월 09:00~11:00", Capacity: 2},
// 		{ID: 2, Name: "Course2", Schedules: "수 14:00~16:00", Capacity: 3},
// 	}
// 	cache, _ := NewEnrollmentCacheWithData(students, courses, nil)

// 	cache.EnrollStudent(1, 1)
// 	cache.AddToWaitlist(2, 1)

// 	info := cache.GetAllCourseCountInfo()

// 	assert.Equal(t, 2, info[1].Capacity)
// 	assert.Equal(t, 1, info[1].EnrolledCount)
// 	assert.Equal(t, 1, info[1].WaitingCount)

// 	assert.Equal(t, 3, info[2].Capacity)
// 	assert.Equal(t, 0, info[2].EnrolledCount)
// 	assert.Equal(t, 0, info[2].WaitingCount)
// }

// func TestLoadEnrollments_WaitlistPosition(t *testing.T) {
// 	students := createTestStudents()
// 	courses := []models.Course{
// 		{ID: 1, Name: "Course1", Schedules: "월 09:00~11:00", Capacity: 1},
// 	}
// 	enrollments := []models.Enrollment{
// 		{StudentID: 1, CourseID: 1, Position: 0, IsWaitlist: false},
// 		{StudentID: 2, CourseID: 1, Position: 0, IsWaitlist: true},
// 		{StudentID: 3, CourseID: 1, Position: 1, IsWaitlist: true},
// 	}

// 	cache, err := NewEnrollmentCacheWithData(students, courses, enrollments)
// 	assert.NoError(t, err)

// 	assert.Equal(t, int32(1), cache.EnrolledCount[1].Load())
// 	assert.Equal(t, int32(2), cache.WaitingCount[1].Load())

// 	_, inCourses := cache.StudentCourses[1][uint(1)]
// 	assert.True(t, inCourses)

// 	_, inWaiting := cache.StudentWaitingCourses[2][uint(1)]
// 	assert.True(t, inWaiting)
// 	_, inWaiting = cache.StudentWaitingCourses[3][uint(1)]
// 	assert.True(t, inWaiting)
// }
