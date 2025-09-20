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

	// 設定 GORM 配置
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	}

	// 連接到 PostgreSQL
	DB, err = gorm.Open(postgres.Open(cfg.Database.DSN), gormConfig)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// 設定連接池
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	// 設定連接池參數
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

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
		// 檢查是否為表已存在的錯誤或其他可忽略的錯誤
		errorStr := fmt.Sprintf("%v", err)
		if strings.Contains(errorStr, "already exists") || 
		   strings.Contains(errorStr, "contains null values") ||
		   strings.Contains(errorStr, "prepared statement") {
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
