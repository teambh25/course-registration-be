package models

type Student struct {
	ID       uint   `gorm:"primaryKey; autoIncrement" json:"-"`
	Name     string `gorm:"not null" json:"name" binding:"required"`
	UserID   string `gorm:"unique; not null" json:"user_id" binding:"required"`
	Password string `gorm:"not null" json:"password" binding:"required"`
}
