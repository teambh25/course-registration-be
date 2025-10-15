package models

type Course struct {
	ID          uint   `gorm:"primaryKey;autoIncrement" json:"-"`
	Name        string `gorm:"unique;not null" json:"name" binding:"required"`
	Instructor  string `gorm:"not null" json:"instructor" binding:"required"`
	Description string `gorm:"type:text;not null" json:"description" binding:"required"` // 소개 (+주의사항)
	Schedules   string `gorm:"type:text;not null" json:"schedules" binding:"required"`   // 강의 시간 (JSON string)
	Capacity    int    `gorm:"not null" json:"capacity" binding:"required"`
	IsSpecial   bool   `gorm:"default:false" json:"is_special"`
}
