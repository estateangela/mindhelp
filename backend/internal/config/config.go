package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config 應用程式配置結構
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	OpenRouter OpenRouterConfig
	CORS     CORSConfig
	Logging  LoggingConfig
}

// ServerConfig 伺服器配置
type ServerConfig struct {
	Port    string
	GinMode string
}

// DatabaseConfig 資料庫配置
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
	DSN      string
}

// JWTConfig JWT 配置
type JWTConfig struct {
	Secret     string
	Expiry     time.Duration
	RefreshExpiry time.Duration
}

// OpenRouterConfig OpenRouter API 配置
type OpenRouterConfig struct {
	APIKey  string
	BaseURL string
}

// CORSConfig CORS 配置
type CORSConfig struct {
	AllowedOrigins []string
}

// LoggingConfig 日誌配置
type LoggingConfig struct {
	Level  string
	Format string
}

// Load 載入配置
func Load() (*Config, error) {
	// 載入 .env 文件
	if err := godotenv.Load(); err != nil {
		// 如果沒有 .env 文件，使用環境變數
		fmt.Println("Warning: .env file not found, using environment variables")
	}

	config := &Config{}

	// 載入伺服器配置
	config.Server = ServerConfig{
		Port:    getEnv("PORT", "8080"),
		GinMode: getEnv("GIN_MODE", "release"),
	}

	// 載入資料庫配置
	config.Database = DatabaseConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "1433"),
		User:     getEnv("DB_USER", "sa"),
		Password: getEnv("DB_PASSWORD", ""),
		Name:     getEnv("DB_NAME", "mindhelp"),
		SSLMode:  getEnv("DB_SSL_MODE", "disable"),
	}

	// 構建 SQL Server DSN
	config.Database.DSN = fmt.Sprintf(
		"sqlserver://%s:%s@%s:%s?database=%s&encrypt=disable",
		config.Database.User,
		config.Database.Password,
		config.Database.Host,
		config.Database.Port,
		config.Database.Name,
	)

	// 載入 JWT 配置
	jwtExpiry, _ := time.ParseDuration(getEnv("JWT_EXPIRY", "24h"))
	refreshExpiry, _ := time.ParseDuration(getEnv("JWT_REFRESH_EXPIRY", "168h")) // 7 days
	
	config.JWT = JWTConfig{
		Secret:     getEnv("JWT_SECRET", "your-super-secret-jwt-key-here"),
		Expiry:     jwtExpiry,
		RefreshExpiry: refreshExpiry,
	}

	// 載入 OpenRouter 配置
	config.OpenRouter = OpenRouterConfig{
		APIKey:  getEnv("OPENROUTER_API_KEY", ""),
		BaseURL: getEnv("OPENROUTER_BASE_URL", "https://openrouter.ai/api/v1"),
	}

	// 載入 CORS 配置
	allowedOrigins := getEnv("ALLOWED_ORIGINS", "http://localhost:3000")
	config.CORS = CORSConfig{
		AllowedOrigins: []string{allowedOrigins},
	}

	// 載入日誌配置
	config.Logging = LoggingConfig{
		Level:  getEnv("LOG_LEVEL", "info"),
		Format: getEnv("LOG_FORMAT", "json"),
	}

	return config, nil
}

// getEnv 獲取環境變數，如果不存在則返回預設值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvInt 獲取整數環境變數
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getEnvBool 獲取布林環境變數
func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}
