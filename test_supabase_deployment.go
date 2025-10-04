package main

import (
	"fmt"
	"os"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	fmt.Println("🔍 Supabase 部署連接測試")
	fmt.Println(strings.Repeat("=", 50))

	// 從環境變數獲取連接資訊
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		// 使用個別參數構建 DSN
		host := getEnv("DB_HOST", "aws-1-ap-southeast-1.pooler.supabase.com")
		port := getEnv("DB_PORT", "6543")
		user := getEnv("DB_USER", "postgres.haunuvdhisdygfradaya")
		password := getEnv("DB_PASSWORD", "MIND_HELP_2025")
		dbname := getEnv("DB_NAME", "postgres")
		sslmode := getEnv("DB_SSL_MODE", "require")

		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s connect_timeout=30",
			host, port, user, password, dbname, sslmode)
	}

	fmt.Printf("📋 測試連接字串: %s\n", maskPassword(dsn))
	fmt.Println()

	// 測試多種配置
	testConfigs := []string{
		dsn, // 當前配置
		strings.Replace(dsn, "sslmode=require", "sslmode=disable", 1), // 無 SSL
		strings.Replace(dsn, "6543", "5432", 1),                       // 標準端口
		strings.Replace(dsn, "pooler.supabase.com", "supabase.co", 1), // 直接連接
	}

	for i, config := range testConfigs {
		fmt.Printf("🧪 測試配置 %d...\n", i+1)
		if testConnection(config) {
			fmt.Printf("✅ 配置 %d 成功!\n", i+1)
			fmt.Printf("📋 建議使用: %s\n", maskPassword(config))
			return
		}
		fmt.Println()
	}

	fmt.Println("❌ 所有配置都失敗了")
	fmt.Println("💡 建議:")
	fmt.Println("1. 檢查 Supabase 專案是否處於活動狀態")
	fmt.Println("2. 確認認證資訊正確")
	fmt.Println("3. 聯繫 Supabase 支援")
}

func testConnection(dsn string) bool {
	// 設定 GORM 配置
	gormConfig := &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	}

	// 嘗試連接
	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		fmt.Printf("❌ 連接失敗: %v\n", err)
		return false
	}

	// 測試 Ping
	sqlDB, err := db.DB()
	if err != nil {
		fmt.Printf("❌ 獲取資料庫實例失敗: %v\n", err)
		return false
	}

	if err := sqlDB.Ping(); err != nil {
		fmt.Printf("❌ Ping 失敗: %v\n", err)
		return false
	}

	// 測試簡單查詢
	var result int
	err = db.Raw("SELECT 1").Scan(&result).Error
	if err != nil {
		fmt.Printf("❌ 查詢測試失敗: %v\n", err)
		return false
	}

	fmt.Printf("✅ 連接成功: %d\n", result)

	// 檢查資料庫統計
	stats := sqlDB.Stats()
	fmt.Printf("📊 連接池狀態: 開啟=%d, 使用中=%d, 閒置=%d\n",
		stats.OpenConnections, stats.InUse, stats.Idle)

	sqlDB.Close()
	return true
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func maskPassword(dsn string) string {
	if len(dsn) > 50 {
		return dsn[:30] + "***" + dsn[len(dsn)-20:]
	}
	return dsn
}
