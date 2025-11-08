package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"course-reg/internal/app/domain/cache"
	"course-reg/internal/app/domain/static"
	"course-reg/internal/app/domain/worker"
	"course-reg/internal/app/handler"
	"course-reg/internal/app/repository"
	"course-reg/internal/app/routers"
	"course-reg/internal/app/service"
	"course-reg/internal/pkg/database"
	"course-reg/internal/pkg/setting"
	"course-reg/internal/pkg/utils"
)

func init() {
	setting.Setup()
	// logging.Setup() // logging 할 수 있는 환경일 떄 다시 사용
	utils.Setup()
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
	registrationStatus := setupRegistrationState()
	services := setupServices(repos, enrollmentWorker, registrationStatus)
	handlers := setupHandlers(services)
	router := setupRouter(handlers, registrationStatus)

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
	return w

	// w.Start(students, courses)
}

func setupRegistrationState() *cache.RegistrationState {
	enabled, startTime, endTime := setting.LoadRegistrationConfig()

	if enabled {
		// todo
	}

	// add registration schedule goroutine

	return cache.NewRegistrationState(enabled, startTime, endTime)
}

func setupServices(repos *Repositories, w *worker.EnrollmentWorker, rc *cache.RegistrationState) *Services {
	return &Services{
		Auth:      service.NewAuthService(repos.Student),
		Admin:     service.NewAdminService(repos.Student, repos.Course, repos.Enrollment, w, rc),
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

func setupRouter(handlers *Handlers, regState *cache.RegistrationState) *gin.Engine {
	return routers.InitRouter(handlers.Admin, handlers.Auth, handlers.CourseReg, regState)
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
