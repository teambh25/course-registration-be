package setting

import (
	"log"
	"sync"
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

type RegistrationStatus struct {
	Enabled   bool
	StartTime string
	EndTime   string
}

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

	ServerSetting.ReadTimeout = ServerSetting.ReadTimeout * time.Second
	ServerSetting.WriteTimeout = ServerSetting.WriteTimeout * time.Second
}

func mapTo(section string, v interface{}) {
	err := cfg.Section(section).MapTo(v)
	if err != nil {
		log.Fatalf("Cfg.MapTo %s err: %v", section, err)
	}
}

// LoadRegistrationConfig loads registration config from ini file
func LoadRegistrationConfig() (enabled bool, startTime, endTime string) {
	var regStatus RegistrationStatus
	mapTo("registration", &regStatus)
	return regStatus.Enabled, regStatus.StartTime, regStatus.EndTime
}

var confWriteMu sync.Mutex

// SaveRegistrationState saves registration enabled state to ini file
func SaveRegistrationState(enabled string) error {
	confWriteMu.Lock()
	defer confWriteMu.Unlock()

	cfg.Section("registration").Key("Enabled").SetValue(enabled)
	err := cfg.SaveTo("conf/app.ini")
	if err != nil {
		return err
	}
	return nil
}

// SaveRegistrationPeriod saves registration period to ini file
func SaveRegistrationPeriod(startTime, endTime string) error {
	confWriteMu.Lock()
	defer confWriteMu.Unlock()

	cfg.Section("registration").Key("StartTime").SetValue(startTime)
	cfg.Section("registration").Key("EndTime").SetValue(endTime)

	err := cfg.SaveTo("conf/app.ini")
	if err != nil {
		return err
	}

	return nil
}
