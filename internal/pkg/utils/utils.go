package utils

import (
	"time"

	"github.com/google/uuid"
)

func GenerateSessionID() string {
	newUUID := uuid.New()
	return newUUID.String()
}

func StringToTime(timeStr string) (time.Time, error) {
	return time.Parse("2006-01-02-15-04", timeStr)
}
