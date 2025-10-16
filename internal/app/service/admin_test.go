package service_test

import (
	"course-reg/internal/app/models"
	"course-reg/internal/app/service"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockStudentRepository struct {
	mock.Mock
}

func (m *MockStudentRepository) GetPassword(username string) (string, error) {
	args := m.Called(username)
	return args.String(0), args.Error(1)
}

func (m *MockStudentRepository) InsertStudents(students []models.Student) error {
	args := m.Called(students)
	return args.Error(0)
}

func (m *MockStudentRepository) DeleteAllStudents() error {
	args := m.Called()
	return args.Error(0)
}

type MockCourseRepository struct {
	mock.Mock
}

func (m *MockCourseRepository) CreateCourse(course *models.Course) error {
	args := m.Called(course)
	course.ID = 1
	return args.Error(0)
}

func (m *MockCourseRepository) DeleteCourse(courseID uint) error {
	args := m.Called(courseID)
	return args.Error(0)
}

type MockEnrollmentRepository struct {
	mock.Mock
}

func setupAdminServiceTest(t *testing.T) (*service.AdminService, *MockStudentRepository, *MockCourseRepository, *MockEnrollmentRepository) {
	mockStudentRepo := new(MockStudentRepository)
	mockCourseRepo := new(MockCourseRepository)
	mockEnrollRepo := new(MockEnrollmentRepository)

	adminService := service.NewAdminService(
		mockStudentRepo,
		mockCourseRepo,
		mockEnrollRepo,
	)

	return adminService, mockStudentRepo, mockCourseRepo, mockEnrollRepo
}

func TestAdminService_RegisterStudents(t *testing.T) {
	adminService, mockStudentRepo, _, _ := setupAdminServiceTest(t)

	students := []models.Student{{Name: "Test Student"}}

	t.Run("Success", func(t *testing.T) {
		mockStudentRepo.On("InsertStudents", students).Return(nil).Once()
		err := adminService.RegisterStudents(students)
		assert.NoError(t, err)
		mockStudentRepo.AssertExpectations(t)
	})

	t.Run("Failure", func(t *testing.T) {
		expectedError := errors.New("insert failed")
		mockStudentRepo.On("InsertStudents", students).Return(expectedError).Once()
		err := adminService.RegisterStudents(students)
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		mockStudentRepo.AssertExpectations(t)
	})
}

func TestAdminService_ResetStudents(t *testing.T) {
	adminService, mockStudentRepo, _, _ := setupAdminServiceTest(t)

	t.Run("Success", func(t *testing.T) {
		mockStudentRepo.On("DeleteAllStudents").Return(nil).Once()
		err := adminService.ResetStudents()
		assert.NoError(t, err)
		mockStudentRepo.AssertExpectations(t)
	})

	t.Run("Failure", func(t *testing.T) {
		expectedError := errors.New("delete failed")
		mockStudentRepo.On("DeleteAllStudents").Return(expectedError).Once()
		err := adminService.ResetStudents()
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		mockStudentRepo.AssertExpectations(t)
	})
}

func TestAdminService_CreateCourse(t *testing.T) {
	adminService, _, mockCourseRepo, _ := setupAdminServiceTest(t)

	course := &models.Course{Name: "Test Course"}

	t.Run("Success", func(t *testing.T) {
		mockCourseRepo.On("CreateCourse", course).Return(nil).Once()
		id, err := adminService.CreateCourse(course)
		assert.NoError(t, err)
		assert.Equal(t, uint(1), id)
		mockCourseRepo.AssertExpectations(t)
	})

	t.Run("Failure", func(t *testing.T) {
		expectedError := errors.New("create failed")
		mockCourseRepo.On("CreateCourse", course).Return(expectedError).Once()
		id, err := adminService.CreateCourse(course)
		assert.Error(t, err)
		assert.Equal(t, uint(0), id)
		assert.Equal(t, expectedError, err)
		mockCourseRepo.AssertExpectations(t)
	})
}

func TestAdminService_DeleteCourse(t *testing.T) {
	adminService, _, mockCourseRepo, _ := setupAdminServiceTest(t)

	courseID := uint(1)

	t.Run("Success", func(t *testing.T) {
		mockCourseRepo.On("DeleteCourse", courseID).Return(nil).Once()
		err := adminService.DeleteCourse(courseID)
		assert.NoError(t, err)
		mockCourseRepo.AssertExpectations(t)
	})

	t.Run("Failure", func(t *testing.T) {
		expectedError := errors.New("delete failed")
		mockCourseRepo.On("DeleteCourse", courseID).Return(expectedError).Once()
		err := adminService.DeleteCourse(courseID)
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		mockCourseRepo.AssertExpectations(t)
	})
}
