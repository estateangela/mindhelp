package database

import (
	"fmt"
	"log"
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

	// 連接到 PostgreSQL - 加入重試機制
	maxRetries := 3
	for i := 0; i < maxRetries; i++ {
		DB, err = gorm.Open(postgres.Open(cfg.Database.DSN), gormConfig)
		if err == nil {
			break
		}
		if i < maxRetries-1 {
			log.Printf("Database connection attempt %d failed, retrying in 5 seconds...", i+1)
			time.Sleep(5 * time.Second)
		}
	}
	if err != nil {
		return fmt.Errorf("failed to connect to database after %d attempts: %w", maxRetries, err)
	}

	// 設定連接池
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	// 設定連接池參數 - 針對 Supabase 優化
	sqlDB.SetMaxIdleConns(5)                   // 減少閒置連接
	sqlDB.SetMaxOpenConns(20)                  // 減少最大連接數
	sqlDB.SetConnMaxLifetime(30 * time.Minute) // 縮短連接生命週期
	sqlDB.SetConnMaxIdleTime(10 * time.Minute) // 設定閒置超時

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
