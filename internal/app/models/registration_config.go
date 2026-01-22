package models

type RegistrationConfig struct {
	ID        uint   `gorm:"primaryKey"`
	Enabled   bool   `gorm:"not null;default:false"`
	StartTime string `gorm:"type:text"`
	EndTime   string `gorm:"type:text"`
}
