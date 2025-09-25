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
	GoogleMaps GoogleMapsConfig
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

// GoogleMapsConfig Google Maps API 配置
type GoogleMapsConfig struct {
	APIKey     string
	BaseURL    string
	GeocodingURL string
	PlacesURL   string
	DirectionsURL string
	DistanceMatrixURL string
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
	if err := godotenv.Load(".env"); err != nil {
		// 嘗試其他可能的路徑
		if err2 := godotenv.Load("../.env"); err2 != nil {
			fmt.Printf("Warning: .env file not found in current or parent directory (%v, %v), using environment variables\n", err, err2)
		}
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
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", ""),
		Name:     getEnv("DB_NAME", "mindhelp"),
		SSLMode:  getEnv("DB_SSL_MODE", "disable"),
	}

	// 構建 PostgreSQL DSN - 支援 Supabase 連接字串格式
	if dsn := getEnv("DATABASE_URL", ""); dsn != "" {
		// 如果提供了完整的 DATABASE_URL，直接使用
		config.Database.DSN = dsn
	} else {
		// 否則使用個別參數構建 - 針對 Supabase 優化
		config.Database.DSN = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s connect_timeout=30",
			config.Database.Host,
			config.Database.Port,
			config.Database.User,
			config.Database.Password,
			config.Database.Name,
			config.Database.SSLMode,
		)
	}

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

	// 載入 Google Maps 配置
	config.GoogleMaps = GoogleMapsConfig{
		APIKey:     getEnv("GOOGLE_MAPS_API_KEY", ""),
		BaseURL:    getEnv("GOOGLE_MAPS_BASE_URL", "https://maps.googleapis.com/maps/api"),
		GeocodingURL: getEnv("GOOGLE_MAPS_GEOCODING_URL", "https://maps.googleapis.com/maps/api/geocode/json"),
		PlacesURL:   getEnv("GOOGLE_MAPS_PLACES_URL", "https://maps.googleapis.com/maps/api/place"),
		DirectionsURL: getEnv("GOOGLE_MAPS_DIRECTIONS_URL", "https://maps.googleapis.com/maps/api/directions/json"),
		DistanceMatrixURL: getEnv("GOOGLE_MAPS_DISTANCE_MATRIX_URL", "https://maps.googleapis.com/maps/api/distancematrix/json"),
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
