package models

import "time"

type Enrollment struct {
	ID         uint      `gorm:"primaryKey;autoIncrement"`
	StudentID  uint      `gorm:"not null;uniqueIndex:idx_student_course"`
	CourseID   uint      `gorm:"not null;uniqueIndex:idx_student_course"`
	Position   int       `gorm:"not null;uniqueIndex:idx_student_course"`
	IsWaitlist bool      `gorm:"default:false"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
}
