package database

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"course-reg/internal/app/models"
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

// WarmupConnectionPool pre-establishes database connections to avoid
// connection creation latency during high-traffic periods (e.g., login surge at registration start)
func WarmupConnectionPool(db *gorm.DB, poolSize int) {
	sqlDB, err := db.DB()
	if err != nil {
		log.Printf("[warn] connection pool warmup failed to get sql.DB: %v", err)
		return
	}

	// todo : ctx input으로 받기
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var wg sync.WaitGroup
	for i := 0; i < poolSize; i++ {
		wg.Go(func() {
			if err := sqlDB.PingContext(ctx); err != nil {
				log.Printf("[warn] connection pool warmup ping failed: %v", err)
			}
		})
	}
	wg.Wait()

	log.Printf("[info] connection pool warmed up with %d connections", poolSize)
}
