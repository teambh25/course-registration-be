package app

import (
	"errors"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"course-reg/internal/app/domain/export"
	"course-reg/internal/app/domain/registration"
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
	Worker   *worker.EnrollmentWorker
	RegState *registration.State
	Router   *gin.Engine
}

const workerQueueSize = 1000

// NewApplication creates and initializes the entire application.
// Dependencies are explicitly passed to each component, so incorrect ordering
// will result in compile errors (undefined variable).
func NewApplication(cfg *setting.Config) (*Application, error) {
	// 1. Database
	db, err := database.Setup(cfg.Database.URL, cfg.Database.PoolSize, cfg.Database.ConnMaxLifetime, cfg.Database.ConnMaxIdleTime)
	if err != nil {
		return nil, fmt.Errorf("database setup failed: %w", err)
	}
	log.Println("[info] database setup completed")

	// 2. Repositories (depends on: db)
	studentRepo := repository.NewStudentRepository(db)
	courseRepo := repository.NewCourseRepository(db)
	enrollRepo := repository.NewEnrollmentRepository(db)
	regConfigRepo := repository.NewRegistrationConfigRepository(db)
	log.Println("[info] repositories setup completed")

	// 3. Static files (depends on: courseRepo)
	if err := export.ExportCoursesToJson(courseRepo); err != nil {
		return nil, fmt.Errorf("static files setup failed: %w", err)
	}
	log.Println("[info] static files setup completed")

	// 4. Worker (depends on: enrollRepo)
	enrollWorker := worker.NewEnrollmentWorker(workerQueueSize, enrollRepo)
	log.Println("[info] worker setup completed")

	// 5. Registration state (depends on: regConfigRepo)
	regState, wasEnabled, err := loadRegistrationState(regConfigRepo)
	if err != nil {
		return nil, fmt.Errorf("registration state setup failed: %w", err)
	}
	log.Printf("[info] registration state setup completed (db_enabled: %v)", wasEnabled)

	// 6. Services (depends on: repos, enrollWorker, regState, db)
	warmup := func() {
		database.WarmupConnectionPool(db, cfg.Database.PoolSize)
	}
	authService := service.NewAuthService(studentRepo, cfg.Secret.AdminID, cfg.Secret.AdminPW)
	adminService := service.NewAdminService(studentRepo, courseRepo, enrollRepo, regConfigRepo, enrollWorker, regState, warmup)
	courseRegService := service.NewCourseRegService(courseRepo, enrollRepo, enrollWorker, regState)
	log.Println("[info] services setup completed")

	// 7. Handlers (depends on: services)
	handlers := &handler.Handlers{
		Auth:      handler.NewAuthHandler(authService),
		Admin:     handler.NewAdminHandler(adminService),
		CourseReg: handler.NewCourseRegHandler(courseRegService),
	}
	log.Println("[info] handlers setup completed")

	// 8. Router (depends on: handlers)
	router := routers.InitRouter(cfg.Server.RunMode, cfg.Secret.SessionKey, handlers)
	log.Println("[info] router setup completed")

	// 9. Restore registration if it was enabled before restart (depends on: adminService)
	if wasEnabled {
		log.Println("[info] restoring registration state from before restart")
		if err := adminService.StartRegistration(); err != nil {
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
				return fmt.Errorf("[warn] failed to close database: %w", err)
			}
		}
	}
	log.Println("[info] application shutdown completed")

	return nil
}

func loadRegistrationState(configRepo repository.RegistrationConfigRepositoryInterface) (*registration.State, bool, error) {
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
	regState := registration.NewState(false, config.StartTime, config.EndTime)
	return regState, config.Enabled, nil
}
