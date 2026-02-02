package database

import (
	"course-reg/internal/app/models"
	"fmt"

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

	// Fix DB connection pool size by setting MaxIdleConns == MaxOpenConns
	sqlDB.SetMaxIdleConns(setting.DatabaseSetting.PoolSize)
	sqlDB.SetMaxOpenConns(setting.DatabaseSetting.PoolSize)

	sqlDB.SetConnMaxLifetime(setting.DatabaseSetting.ConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(setting.DatabaseSetting.ConnMaxIdleTime)
	return db, nil
}
