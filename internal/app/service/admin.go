package service

import (
	"course-reg/internal/app/models"
	"course-reg/internal/app/repository"
	"course-reg/internal/app/worker"
	"course-reg/internal/pkg/setting"
	"log"
)

type AdminService struct {
	studentRepo      repository.StudentRepositoryInterface
	courseRepo       repository.CourseRepositoryInterface
	enrollRepo       repository.EnrollmentRepositoryInterface
	enrollmentWorker *worker.EnrollmentWorker
}

func NewAdminService(
	s repository.StudentRepositoryInterface,
	c repository.CourseRepositoryInterface,
	e repository.EnrollmentRepositoryInterface,
	w *worker.EnrollmentWorker,
) *AdminService {
	return &AdminService{
		studentRepo:      s,
		courseRepo:       c,
		enrollRepo:       e,
		enrollmentWorker: w,
	}
}

func (s *AdminService) RegisterStudents(students []models.Student) error {
	err := s.studentRepo.BulkInsertStudents(students)
	if err != nil {
		log.Println("register students failed:", err.Error())
	}
	return err
}

func (s *AdminService) ResetStudents() error {
	err := s.studentRepo.DeleteAllStudents()
	if err != nil {
		log.Println("reset students failed:", err.Error())
	}
	return err
}

func (s *AdminService) CreateCourse(course *models.Course) (uint, error) {
	err := s.courseRepo.CreateCourse(course)
	if err != nil {
		log.Println("create course failed:", err.Error())
		return 0, err
	}

	// todo : 수강 신청 기간 중에 강의 추가 발생시 (동시성 처리 필요)
	// s.enrollmentWorker.AddCourse(*course)

	return course.ID, nil
}

func (s *AdminService) DeleteCourse(courseID uint) error {
	err := s.courseRepo.DeleteCourse(courseID)
	if err != nil {
		log.Println("delete course failed:", err.Error())
		return err
	}

	// todo : 수강 신청 기간 중에 강의 삭제시 (동시성 처리 필요)
	// s.enrollmentWorker.RemoveCourse(courseID)

	return nil
}

func (s *AdminService) RegisterCourses(courses []models.Course) error {
	// todo : course가 없을 때 처음 한번만 실행 가능하도록

	err := s.courseRepo.BulkInsertCourses(courses)
	if err != nil {
		log.Println("register courses failed:", err.Error())
		return err
	}
	return nil
}

func (s *AdminService) ResetCourses() error {
	// todo : 수강 신청 시작하면 reset 불가능

	err := s.courseRepo.DeleteAllCourses()
	if err != nil {
		log.Println("reset courses failed:", err.Error())
		return err
	}
	return nil
}

func (s *AdminService) SetRegistrationPeriod(startTime, endTime string) error {
	// Validate time format
	_, err := setting.ParsePeriodTime(startTime)
	if err != nil {
		log.Println("invalid start time format:", err.Error())
		return err
	}

	_, err = setting.ParsePeriodTime(endTime)
	if err != nil {
		log.Println("invalid end time format:", err.Error())
		return err
	}

	// Save to config file
	err = setting.SaveRegistrationPeriod(startTime, endTime)
	if err != nil {
		log.Println("failed to save registration period:", err.Error())
		return err
	}

	return nil
}

func (s *AdminService) GetRegistrationPeriod() (string, string, error) {
	startTime := setting.RegistrationPeriodSetting.StartTime
	endTime := setting.RegistrationPeriodSetting.EndTime

	return startTime, endTime, nil
}

// func (s *AdminService) GetEnrolledStudentsByCourse(courseID uint) ([]Student, error)
// func (s *AdminService) ForceEnrollStudent(courseID uint, studentID uint) error
// func (s *AdminService) CancelEnrollment(courseID uint, studentID uint) error
// func (s *AdminService) CheckDuplicateCourses(studentID uint, courseID uint) ([]Course, error)
