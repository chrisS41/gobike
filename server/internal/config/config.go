package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	MongoURI   string
	DBName     string
	JWTSecret  string
	ServerPort string
	LogDir     string
	LogLevel   string
	GinMode    string
}

var cfg *Config

func Load() *Config {
	if err := godotenv.Load("config/gobike.env"); err != nil {
		fmt.Println("Warning: gobike.env 파일을 찾을 수 없습니다. 기본값을 사용합니다.")
	}

	cfg = &Config{
		MongoURI:   getEnv("MONGO_URI", "mongodb://localhost:27017"),
		DBName:     getEnv("DB_NAME", "gobike"),
		JWTSecret:  getEnv("JWT_SECRET", "your-secret-key"),
		ServerPort: getEnv("PORT", "8080"),
		LogDir:     getEnv("LOG_DIR", "logs"),
		LogLevel:   getEnv("LOG_LEVEL", "DEBUG"),
		GinMode:    getEnv("GIN_MODE", "release"),
	}

	if err := validateConfig(cfg); err != nil {
		fmt.Printf("Invalid configuration: %v", err)
	}

	return cfg
}

// validateConfig는 설정값을 검증합니다
func validateConfig(cfg *Config) error {
	if cfg.JWTSecret == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}
	if cfg.MongoURI == "" {
		return fmt.Errorf("MONGO_URI is required")
	}
	return nil
}

// GetConfig는 현재 설정값을 반환합니다
func GetConfig() *Config {
	if cfg == nil {
		return Load()
	}
	return cfg
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
