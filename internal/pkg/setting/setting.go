package setting

import (
	"log"
	"time"

	"github.com/go-ini/ini"
)

type App struct {
	LogSavePath string
	LogSaveName string
	LogFileExt  string
	TimeFormat  string
}

var AppSetting = &App{}

type Server struct {
	RunMode      string
	HttpPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

var ServerSetting = &Server{}

type Secret struct {
	SessionKey string
	AdminID    string
	AdminPW    string
}

var SecretSetting = &Secret{}

type RegistrationPeriod struct {
	StartTime string
	EndTime   string
}

var RegistrationPeriodSetting = &RegistrationPeriod{}

var cfg *ini.File

// Setup initialize the configuration instance
func Setup() {
	var err error
	cfg, err = ini.Load("conf/app.ini")
	if err != nil {
		log.Fatalf("setting.Setup, fail to parse 'conf/app.ini': %v", err)
	}

	mapTo("app", AppSetting)
	mapTo("server", ServerSetting)
	mapTo("secret", SecretSetting)
	mapTo("registration", RegistrationPeriodSetting)

	ServerSetting.ReadTimeout = ServerSetting.ReadTimeout * time.Second
	ServerSetting.WriteTimeout = ServerSetting.WriteTimeout * time.Second
}

// mapTo map section
func mapTo(section string, v interface{}) {
	err := cfg.Section(section).MapTo(v)
	if err != nil {
		log.Fatalf("Cfg.MapTo %s err: %v", section, err)
	}
}

// SaveRegistrationPeriod saves registration period to config file
func SaveRegistrationPeriod(startTime, endTime string) error {
	cfg.Section("registration").Key("StartTime").SetValue(startTime)
	cfg.Section("registration").Key("EndTime").SetValue(endTime)

	err := cfg.SaveTo("conf/app.ini")
	if err != nil {
		return err
	}

	// Update in-memory setting
	RegistrationPeriodSetting.StartTime = startTime
	RegistrationPeriodSetting.EndTime = endTime

	return nil
}

// ParsePeriodTime parses time string in format "yyyy-mm-dd-hh-mm"
func ParsePeriodTime(timeStr string) (time.Time, error) {
	return time.Parse("2006-01-02-15-04", timeStr)
}

// IsWithinRegistrationPeriod checks if given time is within registration period
func IsWithinRegistrationPeriod(now time.Time) (bool, error) {
	if RegistrationPeriodSetting.StartTime == "" || RegistrationPeriodSetting.EndTime == "" {
		return false, nil
	}

	startTime, err := ParsePeriodTime(RegistrationPeriodSetting.StartTime)
	if err != nil {
		return false, err
	}

	endTime, err := ParsePeriodTime(RegistrationPeriodSetting.EndTime)
	if err != nil {
		return false, err
	}

	return now.After(startTime) && now.Before(endTime), nil
}
