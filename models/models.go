package models

// var db *gorm.DB

import (
	"log"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Setup() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("db/course_reg.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database:", err)
		panic(err)
	}

	if err := db.AutoMigrate(&Student{}); err != nil {
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
