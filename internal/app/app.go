package app

import (
	"errors"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"course-reg/internal/app/domain/cache"
	"course-reg/internal/app/domain/export"
	"course-reg/internal/app/domain/worker"
	"course-reg/internal/app/handler"
	"course-reg/internal/app/models"
	"course-reg/internal/app/repository"
	"course-reg/internal/app/routers"
	"course-reg/internal/app/service"
	"course-reg/internal/pkg/database"
	"course-reg/internal/pkg/setting"
)

// Application contains all application components and their dependencies
type Application struct {
	DB       *gorm.DB
	Worker   *worker.Worker
	RegState *cache.RegistrationState
	Router   *gin.Engine
}

const workerQueueSize = 1000

// NewApplication creates and initializes the entire application.
// Dependencies are explicitly passed to each component, so incorrect ordering
// will result in compile errors (undefined variable).
func NewApplication() (*Application, error) {
	// 1. Database
	db, err := database.Setup()
	if err != nil {
		return nil, fmt.Errorf("database setup failed: %w", err)
	}
	log.Println("[info] database setup completed")

	// 2. Repositories (depends on: db)
	repos := newRepositories(db)
	log.Println("[info] repositories setup completed")

	// 3. Static files (depends on: repos.Course)
	if err := export.ExportCoursesToJson(repos.Course); err != nil {
		return nil, fmt.Errorf("static files setup failed: %w", err)
	}
	log.Println("[info] static files setup completed")

	// 4. Worker (depends on: repos.Enrollment)
	enrollWorker := worker.NewEnrollmentWorker(workerQueueSize, repos.Enrollment)
	log.Println("[info] worker setup completed")

	// 5. Registration state (depends on: repos.RegistrationConfig)
	regState, wasEnabled, err := newRegistrationState(repos.RegistrationConfig)
	if err != nil {
		return nil, fmt.Errorf("registration state setup failed: %w", err)
	}
	log.Printf("[info] registration state setup completed (db_enabled: %v)", wasEnabled)

	// 6. Services (depends on: repos, enrollWorker, regState, db)
	warmupFunc := func() {
		database.WarmupConnectionPool(db, setting.DatabaseSetting.PoolSize)
	}
	services := newServices(repos, enrollWorker, regState, warmupFunc)
	log.Println("[info] services setup completed")

	// 7. Handlers (depends on: services)
	handlers := newHandlers(services)
	log.Println("[info] handlers setup completed")

	// 8. Router (depends on: handlers)
	router := routers.InitRouter(handlers.Admin, handlers.Auth, handlers.CourseReg)
	log.Println("[info] router setup completed")

	// 9. Restore registration if it was enabled before restart (depends on: services.Admin)
	if wasEnabled {
		log.Println("[info] restoring registration state from before restart")
		if err := services.Admin.StartRegistration(); err != nil {
			return nil, fmt.Errorf("registration restore failed: %w", err)
		}
	}

	return &Application{
		DB:       db,
		Worker:   enrollWorker,
		RegState: regState,
		Router:   router,
	}, nil
}

// Shutdown gracefully shuts down the application
func (app *Application) Shutdown() error {
	log.Println("[info] shutting down application")

	// Stop worker if running
	if app.Worker != nil && app.RegState != nil && app.RegState.IsEnabled() {
		log.Println("[info] stopping enrollment worker")
		app.Worker.Stop()
	}

	// Close database connection
	if app.DB != nil {
		sqlDB, err := app.DB.DB()
		if err == nil {
			if err := sqlDB.Close(); err != nil {
				return fmt.Errorf("failed to close database: %w", err)
			}
		}
	}

	log.Println("[info] application shutdown completed")
	return nil
}

// --- Helper functions for creating component groups ---

type repositories struct {
	Student            repository.StudentRepositoryInterface
	Course             repository.CourseRepositoryInterface
	Enrollment         repository.EnrollmentRepositoryInterface
	RegistrationConfig repository.RegistrationConfigRepositoryInterface
}

func newRepositories(db *gorm.DB) *repositories {
	return &repositories{
		Student:            repository.NewStudentRepository(db),
		Course:             repository.NewCourseRepository(db),
		Enrollment:         repository.NewEnrollmentRepository(db),
		RegistrationConfig: repository.NewRegistrationConfigRepository(db),
	}
}

func newRegistrationState(configRepo repository.RegistrationConfigRepositoryInterface) (*cache.RegistrationState, bool, error) {
	config, err := configRepo.GetConfig()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			config = &models.RegistrationConfig{
				ID:        1,
				Enabled:   false,
				StartTime: "",
				EndTime:   "",
			}
			if err := configRepo.CreateConfig(config); err != nil {
				return nil, false, fmt.Errorf("failed to create initial registration config: %w", err)
			}
			log.Println("[info] created initial registration config")
		} else {
			return nil, false, fmt.Errorf("failed to load registration config: %w", err)
		}
	}

	// Always start with disabled state; restore later via StartRegistration
	regState := cache.NewRegistrationState(false, config.StartTime, config.EndTime)
	return regState, config.Enabled, nil
}

type services struct {
	Auth      service.AuthServiceInterface
	Admin     service.AdminServiceInterface
	CourseReg service.CourseRegServiceInterface
}

func newServices(repos *repositories, w *worker.Worker, rs *cache.RegistrationState, warmupFunc func()) *services {
	return &services{
		Auth:      service.NewAuthService(repos.Student),
		Admin:     service.NewAdminService(repos.Student, repos.Course, repos.Enrollment, repos.RegistrationConfig, w, rs, warmupFunc),
		CourseReg: service.NewCourseRegService(repos.Course, repos.Enrollment, w, rs),
	}
}

type handlers struct {
	Auth      *handler.AuthHandler
	Admin     *handler.AdminHandler
	CourseReg *handler.CourseRegHandler
}

func newHandlers(s *services) *handlers {
	return &handlers{
		Auth:      handler.NewAuthHandler(s.Auth),
		Admin:     handler.NewAdminHandler(s.Admin),
		CourseReg: handler.NewCourseRegHandler(s.CourseReg),
	}
}
