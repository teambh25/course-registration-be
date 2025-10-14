package models

// var db *gorm.DB

import (
	"database/sql"
	"log"
	"time"

	_ "modernc.org/sqlite" // Pure Go SQLite driver
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Setup() *gorm.DB {
	// Open pure Go SQLite connection
	sqlDB, err := sql.Open("sqlite", "db/course_reg.db")
	if err != nil {
		log.Fatal("failed to open database:", err)
		panic(err)
	}

	// Configure connection pool
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Wrap with GORM
	db, err := gorm.Open(sqlite.Dialector{Conn: sqlDB}, &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database:", err)
		panic(err)
	}

	if err := db.AutoMigrate(&Student{}); err != nil {
		log.Fatal("failed to migrate database:", err)
		panic(err)
	}

	return db
}
