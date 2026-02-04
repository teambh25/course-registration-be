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
)

// Application contains all application components and their dependencies
type Application struct {
	DB       *gorm.DB
	Repos    *Repositories
	Worker   *worker.EnrollmentWorker
	RegState *cache.RegistrationState
	Services *Services
	Handlers *Handlers
	Router   *gin.Engine
}

// Repositories holds all repository instances
type Repositories struct {
	Student            repository.StudentRepositoryInterface
	Course             repository.CourseRepositoryInterface
	Enrollment         repository.EnrollmentRepositoryInterface
	RegistrationConfig repository.RegistrationConfigRepositoryInterface
}

// Services holds all service instances
type Services struct {
	Auth      service.AuthServiceInterface
	Admin     service.AdminServiceInterface
	CourseReg service.CourseRegServiceInterface
}

// Handlers holds all handler instances
type Handlers struct {
	Auth      *handler.AuthHandler
	Admin     *handler.AdminHandler
	CourseReg *handler.CourseRegHandler
}

// NewApplication creates and initializes the entire application
func NewApplication() (*Application, error) {
	app := &Application{}

	if err := app.setupDatabase(); err != nil {
		return nil, err
	}

	if err := app.setupRepositories(); err != nil {
		return nil, err
	}

	if err := app.setupStaticFiles(); err != nil {
		return nil, err
	}

	// Step 4: Setup worker
	if err := app.setupWorker(); err != nil {
		return nil, fmt.Errorf("worker setup failed: %w", err)
	}

	// Step 5: Setup registration state
	if err := app.setupRegistrationState(); err != nil {
		return nil, fmt.Errorf("registration state setup failed: %w", err)
	}

	// Step 6: Load data and start worker if registration is enabled
	if err := app.initializeRegistrationIfEnabled(); err != nil {
		return nil, fmt.Errorf("registration initialization failed: %w", err)
	}

	// Step 7: Setup services
	if err := app.setupServices(); err != nil {
		return nil, fmt.Errorf("service setup failed: %w", err)
	}

	// Step 8: Setup handlers
	if err := app.setupHandlers(); err != nil {
		return nil, fmt.Errorf("handler setup failed: %w", err)
	}

	// Step 9: Setup router
	if err := app.setupRouter(); err != nil {
		return nil, fmt.Errorf("router setup failed: %w", err)
	}

	return app, nil
}

func (app *Application) setupDatabase() error {
	db, err := database.Setup()
	if err != nil {
		return fmt.Errorf("failed to setup DB: %w", err)
	}
	app.DB = db
	log.Println("[info] database setup completed")
	return nil
}

func (app *Application) setupRepositories() error {
	if app.DB == nil {
		return fmt.Errorf("database must be initialized before repositories")
	}

	app.Repos = &Repositories{
		Student:            repository.NewStudentRepository(app.DB),
		Course:             repository.NewCourseRepository(app.DB),
		Enrollment:         repository.NewEnrollmentRepository(app.DB),
		RegistrationConfig: repository.NewRegistrationConfigRepository(app.DB),
	}
	log.Println("[info] repositories setup completed")
	return nil
}

func (app *Application) setupStaticFiles() error {
	if err := export.ExportCoursesToJson(app.Repos.Course); err != nil {
		return fmt.Errorf("failed to export courses to JSON: %w", err)
	}
	log.Println("[info] static files setup completed")
	return nil
}

// setupWorker initializes the enrollment worker
func (app *Application) setupWorker() error {
	// todo: read queue size from setting variables
	const defaultQueueSize = 1000
	app.Worker = worker.NewEnrollmentWorker(defaultQueueSize, app.Repos.Enrollment)
	log.Println("[info] enrollment worker setup completed")
	return nil
}

// setupRegistrationState initializes the registration state
func (app *Application) setupRegistrationState() error {
	config, err := app.Repos.RegistrationConfig.GetConfig()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 초기 레코드가 없으면 기본값으로 생성
			config = &models.RegistrationConfig{
				ID:        1,
				Enabled:   false,
				StartTime: "",
				EndTime:   "",
			}
			if err := app.Repos.RegistrationConfig.CreateConfig(config); err != nil {
				return fmt.Errorf("failed to create initial registration config: %w", err)
			}
			log.Println("[info] created initial registration config")
		} else {
			return fmt.Errorf("failed to load registration config: %w", err)
		}
	}
	app.RegState = cache.NewRegistrationState(config.Enabled, config.StartTime, config.EndTime)
	log.Printf("[info] registration state setup completed (enabled: %v, start: %s, end: %s)",
		config.Enabled, config.StartTime, config.EndTime)
	return nil
}

// initializeRegistrationIfEnabled loads data and starts worker if registration is enabled
func (app *Application) initializeRegistrationIfEnabled() error {
	if !app.RegState.IsEnabled() {
		log.Println("[info] registration is disabled, skipping worker initialization")
		return nil
	}

	log.Println("[info] registration is enabled, loading data and starting worker")

	// Load students
	students, err := app.Repos.Student.FetchAllStudents()
	if err != nil {
		return fmt.Errorf("failed to load students: %w", err)
	}

	// Load courses
	courses, err := app.Repos.Course.FetchAllCourses()
	if err != nil {
		return fmt.Errorf("failed to load courses: %w", err)
	}

	// Load enrollments
	enrollments, err := app.Repos.Enrollment.FetchAllEnrollments()
	if err != nil {
		return fmt.Errorf("failed to load enrollments: %w", err)
	}

	// Start worker with loaded data
	if err := app.Worker.Start(students, courses, enrollments); err != nil {
		return fmt.Errorf("failed to start worker: %w", err)
	}

	log.Println("[info] worker started with loaded data")
	return nil
}

// setupServices initializes all services
func (app *Application) setupServices() error {
	app.Services = &Services{
		Auth:      service.NewAuthService(app.Repos.Student),
		Admin:     service.NewAdminService(app.Repos.Student, app.Repos.Course, app.Repos.Enrollment, app.Repos.RegistrationConfig, app.Worker, app.RegState),
		CourseReg: service.NewCourseRegService(app.Repos.Course, app.Repos.Enrollment, app.Worker, app.RegState),
	}
	log.Println("[info] services setup completed")
	return nil
}

// setupHandlers initializes all handlers
func (app *Application) setupHandlers() error {
	app.Handlers = &Handlers{
		Auth:      handler.NewAuthHandler(app.Services.Auth),
		Admin:     handler.NewAdminHandler(app.Services.Admin),
		CourseReg: handler.NewCourseRegHandler(app.Services.CourseReg),
	}
	log.Println("[info] handlers setup completed")
	return nil
}

// setupRouter initializes the HTTP router
func (app *Application) setupRouter() error {
	app.Router = routers.InitRouter(app.Handlers.Admin, app.Handlers.Auth, app.Handlers.CourseReg)
	log.Println("[info] router setup completed")
	return nil
}

// Shutdown gracefully shuts down the application
func (app *Application) Shutdown() error {
	log.Println("[info] shutting down application")

	// Stop worker if running
	if app.Worker != nil && app.RegState != nil && app.RegState.IsEnabled() {
		log.Println("[info] stopping enrollment worker")
		app.Worker.Stop()
	}

	// Close database connection if needed
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
