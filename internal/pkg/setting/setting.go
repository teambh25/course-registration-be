package setting

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	App      App
	Server   Server
	Secret   Secret
	Database Database
}

type App struct {
	LogSavePath string
	LogSaveName string
	LogFileExt  string
	TimeFormat  string
}

type Server struct {
	RunMode      string
	HttpPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type Secret struct {
	SessionKey string
	AdminID    string
	AdminPW    string
}

type Database struct {
	URL             string
	PoolSize        int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

// Load reads environment variables and returns a Config instance
func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		App: App{
			LogSavePath: getEnvRequired("APP_LOG_SAVE_PATH"),
			LogSaveName: getEnvRequired("APP_LOG_SAVE_NAME"),
			LogFileExt:  getEnvRequired("APP_LOG_FILE_EXT"),
			TimeFormat:  getEnvRequired("APP_TIME_FORMAT"),
		},
		Server: Server{
			RunMode:      getEnvRequired("SERVER_RUN_MODE"),
			HttpPort:     getEnvAsIntRequired("SERVER_HTTP_PORT"),
			ReadTimeout:  time.Duration(getEnvAsIntRequired("SERVER_READ_TIMEOUT")) * time.Second,
			WriteTimeout: time.Duration(getEnvAsIntRequired("SERVER_WRITE_TIMEOUT")) * time.Second,
		},
		Database: Database{
			URL:             getEnvRequired("DATABASE_URL"),
			PoolSize:        getEnvAsIntRequired("DATABASE_POOL_SIZE"),
			ConnMaxLifetime: time.Duration(getEnvAsIntRequired("DATABASE_CONN_MAX_LIFETIME")) * time.Minute,
			ConnMaxIdleTime: time.Duration(getEnvAsIntRequired("DATABASE_CONN_MAX_IDLE_TIME")) * time.Minute,
		},
		Secret: Secret{
			SessionKey: getEnvRequired("SECRET_SESSION_KEY"),
			AdminID:    getEnvRequired("SECRET_ADMIN_ID"),
			AdminPW:    getEnvRequired("SECRET_ADMIN_PW"),
		},
	}
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
