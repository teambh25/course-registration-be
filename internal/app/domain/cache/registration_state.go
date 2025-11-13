package cache

import (
	"course-reg/internal/pkg/utils"
	"fmt"
	"sync"
	"time"
)

type RegistrationState struct {
	mu        sync.RWMutex
	enabled   bool
	startTime string
	endTime   string
}

func NewRegistrationState(enabled bool, startTime, endTime string) *RegistrationState {
	return &RegistrationState{
		enabled:   enabled,
		startTime: startTime,
		endTime:   endTime,
	}
}

func (rs *RegistrationState) IsEnabled() bool {
	rs.mu.RLock()
	defer rs.mu.RUnlock()
	return rs.enabled
}

func (rs *RegistrationState) ChangeEnabledAndAct(enabled bool, act func() error) error {
	rs.mu.Lock()
	defer rs.mu.Unlock()

	if rs.enabled == enabled {
		return fmt.Errorf("enabled aleardy set %v", enabled)
	}

	if err := act(); err != nil {
		return err
	}

	rs.enabled = enabled
	return nil
}

func (rs *RegistrationState) RunIfEnabled(enabled bool, act func() error) error {
	if rs.mu.TryRLock() {
		defer rs.mu.RUnlock()
		if rs.enabled != enabled {
			return fmt.Errorf("enabled is not %v", enabled)
		}
		if err := act(); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("registration is setting up")
	}
	return nil
}

// GetPeriod returns the registration start and end times
func (rs *RegistrationState) GetPeriod() (startTime, endTime string) {
	rs.mu.RLock()
	defer rs.mu.RUnlock()
	return rs.startTime, rs.endTime
}

// SetPeriod sets the registration period
func (rs *RegistrationState) SetPeriod(startTime, endTime string) {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	rs.startTime = startTime
	rs.endTime = endTime
}

// IsWithinRegistrationPeriod checks if the given time is within the registration period
func (rs *RegistrationState) IsWithinRegistrationPeriod(now time.Time) (bool, error) {
	rs.mu.RLock()
	defer rs.mu.RUnlock()

	if rs.startTime == "" || rs.endTime == "" {
		return false, nil
	}

	startTime, err := utils.StringToTime(rs.startTime)
	if err != nil {
		return false, err
	}

	endTime, err := utils.StringToTime(rs.endTime)
	if err != nil {
		return false, err
	}

	return now.After(startTime) && now.Before(endTime), nil
}
