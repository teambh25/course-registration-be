package setting

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
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

type Database struct {
	URL          string
	MaxConns     int
	MaxIdleConns int
}

var DatabaseSetting = &Database{}

// Setup initialize the configuration instance from environment variables
func Setup() {
	// Load .env file if it exists (ignore error if file doesn't exist)
	_ = godotenv.Load()

	// App settings
	AppSetting.LogSavePath = getEnvRequired("APP_LOG_SAVE_PATH")
	AppSetting.LogSaveName = getEnvRequired("APP_LOG_SAVE_NAME")
	AppSetting.LogFileExt = getEnvRequired("APP_LOG_FILE_EXT")
	AppSetting.TimeFormat = getEnvRequired("APP_TIME_FORMAT")

	// Server settings
	ServerSetting.RunMode = getEnvRequired("SERVER_RUN_MODE")
	ServerSetting.HttpPort = getEnvAsIntRequired("SERVER_HTTP_PORT")
	ServerSetting.ReadTimeout = time.Duration(getEnvAsIntRequired("SERVER_READ_TIMEOUT")) * time.Second
	ServerSetting.WriteTimeout = time.Duration(getEnvAsIntRequired("SERVER_WRITE_TIMEOUT")) * time.Second

	// Database settings
	DatabaseSetting.URL = getEnvRequired("DATABASE_URL")
	DatabaseSetting.MaxConns = getEnvAsIntRequired("DATABASE_MAX_CONNS")
	DatabaseSetting.MaxIdleConns = getEnvAsIntRequired("DATABASE_MAX_IDLE_CONNS")

	// Secret settings
	SecretSetting.SessionKey = getEnvRequired("SECRET_SESSION_KEY")
	SecretSetting.AdminID = getEnvRequired("SECRET_ADMIN_ID")
	SecretSetting.AdminPW = getEnvRequired("SECRET_ADMIN_PW")
}

// getEnvRequired retrieves an environment variable or exits if not set
func getEnvRequired(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Required environment variable %s is not set", key)
	}
	return value
}

// getEnvAsIntRequired retrieves an environment variable as an integer or exits if not set/invalid
func getEnvAsIntRequired(key string) int {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Required environment variable %s is not set", key)
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		log.Fatalf("Environment variable %s must be a valid integer, got: %s", key, value)
	}
	return intValue
}
