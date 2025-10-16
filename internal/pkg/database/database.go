package database

import (
	"course-reg/internal/app/models"
	"log"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Setup() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("db/course_reg.db"), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		log.Fatal("failed to connect database:", err)
		panic(err)
	}

	if err := db.AutoMigrate(&models.Student{}, &models.Course{}); err != nil {
		log.Fatal("failed to migrate database:", err)
		panic(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("failed to return sql.DB:", err)
		panic(err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db
}
