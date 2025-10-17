package models

import "time"

type Enrollment struct {
	ID               uint      `gorm:"primaryKey;autoIncrement"`
	StudentID        uint      `gorm:"not null;index:idx_student_course,unique"`
	CourseID         uint      `gorm:"not null;index:idx_student_course,unique"`
	IsWaitlist       bool      `gorm:"default:false"`
	WaitlistPosition int       `gorm:"default:0"`
	CreatedAt        time.Time `gorm:"autoCreateTime"`
}
