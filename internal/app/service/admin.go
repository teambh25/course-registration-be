package service

import (
	"course-reg/internal/app/domain/cache"
	"course-reg/internal/app/domain/export"
	"course-reg/internal/app/domain/worker"
	"course-reg/internal/app/models"
	"course-reg/internal/app/repository"
	"course-reg/internal/pkg/setting"
	"course-reg/internal/pkg/utils"
	"log"
)

type AdminService struct {
	studentRepo  repository.StudentRepositoryInterface
	courseRepo   repository.CourseRepositoryInterface
	enrollRepo   repository.EnrollmentRepositoryInterface
	enrollWorker *worker.EnrollmentWorker
	regState     *cache.RegistrationState
}

func NewAdminService(
	s repository.StudentRepositoryInterface,
	c repository.CourseRepositoryInterface,
	e repository.EnrollmentRepositoryInterface,
	w *worker.EnrollmentWorker,
	rc *cache.RegistrationState,
) *AdminService {
	return &AdminService{
		studentRepo:  s,
		courseRepo:   c,
		enrollRepo:   e,
		enrollWorker: w,
		regState:     rc,
	}
}

func (s *AdminService) GetRegistrationState() bool {
	return s.regState.IsEnabled()
}

func (s *AdminService) StartRegistration() error {
	err := s.regState.ChangeEnabledAndAct(true, func() error {

		students, err := s.studentRepo.FetchAllStudents()
		if err != nil {
			log.Println("failed to load students:", err.Error())
			return err
		}

		courses, err := s.courseRepo.FetchAllCourses()
		if err != nil {
			log.Println("failed to load courses:", err.Error())
			return err
		}

		enrollments, err := s.enrollRepo.FetchAllEnrollments()
		if err != nil {
			log.Printf("failed to load enrollments: %v", err)
			return err
		}

		if err := s.enrollWorker.Start(students, courses, enrollments); err != nil {
			log.Printf("failed to start worker: %v", err)
			return err
		}

		if err := setting.SaveRegistrationState("true"); err != nil {
			log.Println("save registration state failed:", err.Error())
			return err
		}

		return nil
	})

	if err != nil {
		log.Println(err)
	} else {
		log.Println("Start Registration!!!")
	}

	return err
}

func (s *AdminService) PauseRegistration() error {
	err := s.regState.ChangeEnabledAndAct(false, func() error {
		s.enrollWorker.Stop()

		if err := setting.SaveRegistrationState("false"); err != nil {
			log.Println("save registration state failed:", err.Error()) // 500
			return err
		}
		return nil
	})
	if err != nil {
		log.Println(err)
	} else {
		log.Println("Pause Registration!!!")
	}
	return err
}

func (s *AdminService) GetRegistrationPeriod() (string, string) {
	return s.regState.GetPeriod()
}

func (s *AdminService) SetRegistrationPeriod(startTime, endTime string) error {
	// Validate time format
	_, err := utils.StringToTime(startTime)
	if err != nil {
		log.Println("invalid start time format:", err.Error())
		return err
	}

	_, err = utils.StringToTime(endTime)
	if err != nil {
		log.Println("invalid end time format:", err.Error())
		return err
	}

	// Update in-memory state
	s.regState.SetPeriod(startTime, endTime)

	// Save to config file
	err = setting.SaveRegistrationPeriod(startTime, endTime)
	if err != nil {
		log.Println("failed to save registration period:", err.Error())
		return err
	}

	return nil
}

func (s *AdminService) RegisterStudents(students []models.Student) error {
	err := s.regState.RunIfEnabled(false, func() error {
		return s.studentRepo.BulkInsertStudents(students)
	})
	if err != nil {
		log.Println("register students failed:", err.Error())
		return err
	}
	return nil
}

func (s *AdminService) ResetStudents() error {
	err := s.regState.RunIfEnabled(false, func() error {
		return s.studentRepo.DeleteAllStudents()
	})
	if err != nil {
		log.Println("reset students failed:", err.Error())
		return err
	}

	return nil
}

func (s *AdminService) CreateCourse(course *models.Course) (uint, error) {
	err := s.regState.RunIfEnabled(false, func() error {
		return s.courseRepo.InsertCourse(course)
	})
	if err != nil {
		log.Println("create course failed:", err.Error())
		return 0, err
	}

	export.ExportCoursesToJson(s.courseRepo)
	return course.ID, nil
}

func (s *AdminService) DeleteCourse(courseID uint) error {
	err := s.regState.RunIfEnabled(false, func() error {
		return s.courseRepo.DeleteCourse(courseID)
	})
	if err != nil {
		log.Println("delete course failed:", err.Error())
		return err
	}

	export.ExportCoursesToJson(s.courseRepo)
	return nil
}

func (s *AdminService) RegisterCourses(courses []models.Course) error {
	// todo: course가 없을 때만 실행 가능하도록?
	// todo: shcedule에 대한 validation?

	err := s.regState.RunIfEnabled(false, func() error {
		return s.courseRepo.BulkInsertCourses(courses)
	})
	if err != nil {
		log.Println("register courses failed:", err.Error())
		return err
	}

	export.ExportCoursesToJson(s.courseRepo)
	return nil
}

func (s *AdminService) ResetCourses() error {
	err := s.regState.RunIfEnabled(false, func() error {
		return s.courseRepo.DeleteAllCourses()
	})
	if err != nil {
		log.Println("reset courses failed:", err.Error())
		return err
	}

	export.ExportCoursesToJson(s.courseRepo)
	return nil
}

// func (s *AdminService) GetEnrolledStudentsByCourse(courseID uint) ([]Student, error)
// func (s *AdminService) ForceEnrollStudent(courseID uint, studentID uint) error
// func (s *AdminService) CancelEnrollment(courseID uint, studentID uint) error
// func (s *AdminService) CheckDuplicateCourses(studentID uint, courseID uint) ([]Course, error)
