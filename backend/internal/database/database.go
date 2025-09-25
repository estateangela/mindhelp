package database

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"mindhelp-backend/internal/config"
	"mindhelp-backend/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB 全域資料庫實例
var DB *gorm.DB

// Connect 連接到資料庫
func Connect(cfg *config.Config) error {
	var err error

	// 設定 GORM 配置 - 針對 Supabase 優化
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn), // 減少日誌輸出
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
		PrepareStmt:                              false, // 禁用 prepared statements 避免重複錯誤
		DisableForeignKeyConstraintWhenMigrating: true,  // 禁用外鍵約束
	}

	// 連接到 PostgreSQL - 加入重試機制和 Supabase 優化
	maxRetries := 10              // 增加重試次數
	retryDelay := 5 * time.Second // 增加初始延遲

	for i := 0; i < maxRetries; i++ {
		log.Printf("嘗試連接資料庫 (第 %d/%d 次)...", i+1, maxRetries)

		// 嘗試重新構建 DSN 以確保最新配置
		if dsn := getEnv("DATABASE_URL", ""); dsn != "" {
			// 如果提供了完整的 DATABASE_URL，直接使用
			cfg.Database.DSN = dsn
		} else {
			// 否則使用個別參數構建
			cfg.Database.DSN = fmt.Sprintf(
				"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s connect_timeout=30",
				cfg.Database.Host,
				cfg.Database.Port,
				cfg.Database.User,
				cfg.Database.Password,
				cfg.Database.Name,
				cfg.Database.SSLMode,
			)
		}

		DB, err = gorm.Open(postgres.Open(cfg.Database.DSN), gormConfig)
		if err == nil {
			// 測試連接是否真的可用
			sqlDB, testErr := DB.DB()
			if testErr == nil {
				if pingErr := sqlDB.Ping(); pingErr == nil {
					log.Printf("資料庫連接成功!")
					break
				} else {
					log.Printf("資料庫 Ping 失敗: %v", pingErr)
					err = pingErr
				}
			} else {
				log.Printf("無法獲取資料庫實例: %v", testErr)
				err = testErr
			}
		}

		if err != nil {
			log.Printf("資料庫連接失敗 (第 %d/%d 次): %v", i+1, maxRetries, err)
			if i < maxRetries-1 {
				log.Printf("等待 %v 後重試...", retryDelay)
				time.Sleep(retryDelay)
				retryDelay = time.Duration(float64(retryDelay) * 1.5) // 指數退避
				// 最多等待 60 秒
				if retryDelay > 60*time.Second {
					retryDelay = 60 * time.Second
				}
			}
		}
	}

	if err != nil {
		return fmt.Errorf("資料庫連接失敗，已重試 %d 次: %w", maxRetries, err)
	}

	// 設定連接池
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	// 設定連接池參數 - 針對 Supabase 優化
	sqlDB.SetMaxIdleConns(2)                   // Supabase 限制較嚴，減少閒置連接
	sqlDB.SetMaxOpenConns(10)                  // Supabase 限制較嚴，減少最大連接數
	sqlDB.SetConnMaxLifetime(15 * time.Minute) // Supabase 連接超時較短
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)  // 縮短閒置超時

	log.Println("Successfully connected to database")
	return nil
}

// Migrate 執行資料庫遷移
func Migrate() error {
	log.Println("Starting database migration...")

	// 自動遷移所有模型 - AutoMigrate 會自動處理已存在的表
	err := DB.AutoMigrate(
		&models.User{},
		&models.ChatMessage{},
		&models.ChatSession{},
		&models.Location{},
		&models.Article{},
		&models.Quiz{},
		&models.Review{},
		&models.Bookmark{},
		&models.Notification{},
		&models.UserSetting{},
		&models.AppConfig{},
		&models.Share{},
		&models.Counselor{},
		&models.CounselingCenter{},
		&models.RecommendedDoctor{},
	)
	if err != nil {
		// 檢查是否為可忽略的錯誤
		errorStr := fmt.Sprintf("%v", err)
		if strings.Contains(errorStr, "already exists") ||
			strings.Contains(errorStr, "contains null values") ||
			strings.Contains(errorStr, "prepared statement") ||
			strings.Contains(errorStr, "duplicate key") ||
			strings.Contains(errorStr, "constraint") {
			log.Printf("Migration warning (ignoring): %v", err)
			return nil
		}
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Println("Database migration completed successfully")
	return nil
}

// Close 關閉資料庫連接
func Close() error {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			return fmt.Errorf("failed to get database instance: %w", err)
		}
		return sqlDB.Close()
	}
	return nil
}

// getEnv 獲取環境變數，如果不存在則返回預設值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetDB 獲取資料庫實例
func GetDB() *gorm.DB {
	return DB
}

// CheckConnection 檢查資料庫連接狀態
func CheckConnection() error {
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	return nil
}

// GetConnectionStats 獲取連接池統計資訊
func GetConnectionStats() (map[string]interface{}, error) {
	if DB == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	stats := sqlDB.Stats()
	return map[string]interface{}{
		"open_connections":    stats.OpenConnections,
		"in_use":              stats.InUse,
		"idle":                stats.Idle,
		"wait_count":          stats.WaitCount,
		"wait_duration":       stats.WaitDuration.String(),
		"max_idle_closed":     stats.MaxIdleClosed,
		"max_lifetime_closed": stats.MaxLifetimeClosed,
	}, nil
}

// Reconnect 重新連接資料庫
func Reconnect(cfg *config.Config) error {
	log.Println("嘗試重新連接資料庫...")

	// 關閉現有連接
	if DB != nil {
		if sqlDB, err := DB.DB(); err == nil {
			sqlDB.Close()
		}
	}

	// 重新連接
	return Connect(cfg)
}

// IsHealthy 檢查資料庫健康狀態
func IsHealthy() bool {
	if DB == nil {
		return false
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return false
	}

	if err := sqlDB.Ping(); err != nil {
		return false
	}

	return true
}
