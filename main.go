package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"course-reg/internal/app/domain/static"
	"course-reg/internal/app/domain/worker"
	"course-reg/internal/app/handler"
	"course-reg/internal/app/repository"
	"course-reg/internal/app/routers"
	"course-reg/internal/app/service"
	"course-reg/internal/pkg/database"
	"course-reg/internal/pkg/setting"
	"course-reg/internal/pkg/util"
)

func init() {
	setting.Setup()
	// logging.Setup() // logging 할 수 있는 환경일 떄 다시 사용
	util.Setup()
}

type Repositories struct {
	Student    repository.StudentRepositoryInterface
	Course     repository.CourseRepositoryInterface
	Enrollment repository.EnrollmentRepositoryInterface
}

type Services struct {
	Auth      service.AuthServiceInterface
	Admin     service.AdminServiceInterface
	CourseReg service.CourseRegServiceInterface
}

type Handlers struct {
	Auth      *handler.AuthHandler
	Admin     *handler.AdminHandler
	CourseReg *handler.CourseRegHandler
}

func main() {
	db := database.Setup()
	repos := setupRepositories(db)
	setupStatic(repos)
	enrollmentWorker := setupWorker(repos)
	services := setupServices(repos, enrollmentWorker)
	handlers := setupHandlers(services)
	router := setupRouter(handlers)

	startServer(router)
}

func setupRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		Student:    repository.NewStudentRepository(db),
		Course:     repository.NewCourseRepository(db),
		Enrollment: repository.NewEnrollmentRepository(db),
	}
}

func setupStatic(repos *Repositories) {
	err := static.ExportCoursesToJson(repos.Course)
	if err != nil {
		log.Fatalf("failed to export course json: %v", err)
	}
}

func setupWorker(repos *Repositories) *worker.EnrollmentWorker {
	w := worker.NewEnrollmentWorker(1000)

	courses, err := repos.Course.FetchAllCourses()
	if err != nil {
		log.Fatalf("failed to load courses: %v", err)
	}

	students, err := repos.Student.FetchAllStudents()
	if err != nil {
		log.Fatalf("failed to load students: %v", err)
	}

	// TODO: Load previous enrollment data
	// enrollments, err := repos.Enrollment.LoadAllEnrollments()
	// if err != nil {
	// 	log.Printf("[warning] failed to load enrollments: %v", err)
	// 	enrollments = []models.Enrollment{}
	// }

	w.Start(students, courses)

	return w
}

func setupServices(repos *Repositories, w *worker.EnrollmentWorker) *Services {
	return &Services{
		Auth:      service.NewAuthService(repos.Student),
		Admin:     service.NewAdminService(repos.Student, repos.Course, repos.Enrollment, w),
		CourseReg: service.NewCourseRegService(repos.Course, repos.Enrollment, w),
	}
}

func setupHandlers(services *Services) *Handlers {
	return &Handlers{
		Auth:      handler.NewAuthHandler(services.Auth),
		Admin:     handler.NewAdminHandler(services.Admin),
		CourseReg: handler.NewCourseRegHandler(services.CourseReg),
	}
}

func setupRouter(handlers *Handlers) *gin.Engine {
	timeProvider := util.NewKoreaTimeProvider()
	return routers.InitRouter(handlers.Admin, handlers.Auth, handlers.CourseReg, timeProvider)
}

func startServer(router *gin.Engine) {
	server := &http.Server{
		Addr:           fmt.Sprintf(":%d", setting.ServerSetting.HttpPort),
		Handler:        router,
		ReadTimeout:    setting.ServerSetting.ReadTimeout,
		WriteTimeout:   setting.ServerSetting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	log.Printf("[info] start http server listening %s", server.Addr)

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
