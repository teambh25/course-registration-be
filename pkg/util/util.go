package util

import (
	"github.com/google/uuid"
)

// Setup Initialize the util
func Setup() {
	// jwtSecret = []byte(setting.AppSetting.JwtSecret)
}

func GenerateSessionID() string {
	newUUID := uuid.New()
	return newUUID.String()
}
