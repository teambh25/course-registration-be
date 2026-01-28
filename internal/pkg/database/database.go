package database

import (
	"course-reg/internal/app/models"
	"fmt"
	"time"

	"course-reg/internal/pkg/setting"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Setup() (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(setting.DatabaseSetting.URL), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	if err := db.AutoMigrate(&models.Student{}, &models.Course{}, &models.Enrollment{}, &models.RegistrationConfig{}); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to return sql.DB: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}
