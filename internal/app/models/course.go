package models

type Course struct {
	ID          uint   `gorm:"primaryKey;autoIncrement"` // `gorm:"primaryKey;autoIncrement" json:"-"`
	Name        string `gorm:"unique;not null" json:"name" binding:"required"`
	Instructor  string `gorm:"not null" json:"instructor" binding:"required"`
	Description string `gorm:"type:text" json:"description"`
	Schedules   string `gorm:"type:text;not null" json:"schedules" binding:"required"`
	Capacity    int    `gorm:"not null" json:"capacity" binding:"required"`
	IsSpecial   bool   `gorm:"default:false" json:"is_special"`
}
