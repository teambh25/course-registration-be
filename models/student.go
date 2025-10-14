package models

type Student struct {
	ID          uint   `gorm:"primaryKey; autoIncrement" json:"-"`
	Name        string `gorm:"not null" json:"name" binding:"required"`
	BirthDate   string `gorm:"not null" json:"birth_date" binding:"required"`
	PhoneNumber string `gorm:"unique; not null" json:"phone_number" binding:"required"`
}
