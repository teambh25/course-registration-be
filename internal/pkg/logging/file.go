package logging

import (
	"fmt"
	"time"

	"course-reg/internal/pkg/setting"
)

// getLogFilePath get the log file save path
func getLogFilePath() string {
	return setting.AppSetting.LogSavePath
}

// getLogFileName get the save name of the log file
func getLogFileName() string {
	return fmt.Sprintf("%s%s.%s",
		setting.AppSetting.LogSaveName,
		time.Now().Format(setting.AppSetting.TimeFormat),
		setting.AppSetting.LogFileExt,
	)
}
