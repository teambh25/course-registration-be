package repository

import (
	"course-reg/internal/app/models"

	"gorm.io/gorm"
)

const defaultConfigID = 1

type RegistrationConfigRepository struct {
	db *gorm.DB
}

func NewRegistrationConfigRepository(db *gorm.DB) *RegistrationConfigRepository {
	return &RegistrationConfigRepository{db: db}
}

func (r *RegistrationConfigRepository) GetConfig() (*models.RegistrationConfig, error) {
	var config models.RegistrationConfig
	if err := r.db.First(&config, defaultConfigID).Error; err != nil {
		return nil, err
	}
	return &config, nil
}

func (r *RegistrationConfigRepository) CreateConfig(config *models.RegistrationConfig) error {
	return r.db.Create(config).Error
}

func (r *RegistrationConfigRepository) UpdateEnabled(enabled bool) error {
	return r.db.Model(&models.RegistrationConfig{}).
		Where("id = ?", defaultConfigID).
		Update("enabled", enabled).Error
}

func (r *RegistrationConfigRepository) UpdatePeriod(startTime, endTime string) error {
	return r.db.Model(&models.RegistrationConfig{}).
		Where("id = ?", defaultConfigID).
		Updates(map[string]interface{}{
			"start_time": startTime,
			"end_time":   endTime,
		}).Error
}
