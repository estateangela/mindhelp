package main

import (
	"database/sql"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	fmt.Println("🔍 Supabase 連接測試")
	fmt.Println("=" + strings.Repeat("=", 50))

	// 從環境變數或直接設定連接字串
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

	fmt.Printf("📋 連接字串: %s\n", maskPassword(dsn))
	fmt.Println()

	// 測試連接
	fmt.Println("🧪 測試資料庫連接...")

	maxRetries := 5
	retryDelay := 2 * time.Second

	for i := 0; i < maxRetries; i++ {
		fmt.Printf("嘗試連接 (第 %d/%d 次)...\n", i+1, maxRetries)

		db, err := sql.Open("postgres", dsn)
		if err != nil {
			fmt.Printf("❌ 開啟連接失敗: %v\n", err)
			if i < maxRetries-1 {
				fmt.Printf("等待 %v 後重試...\n", retryDelay)
				time.Sleep(retryDelay)
				retryDelay *= 2
			}
			continue
		}

		// 測試 Ping
		if err := db.Ping(); err != nil {
			fmt.Printf("❌ Ping 失敗: %v\n", err)
			db.Close()
			if i < maxRetries-1 {
				fmt.Printf("等待 %v 後重試...\n", retryDelay)
				time.Sleep(retryDelay)
				retryDelay *= 2
			}
			continue
		}

		fmt.Printf("✅ 連接成功!\n")

		// 測試查詢
		fmt.Println("\n🧪 測試查詢...")
		var result int
		err = db.QueryRow("SELECT 1").Scan(&result)
		if err != nil {
			fmt.Printf("❌ 查詢失敗: %v\n", err)
		} else {
			fmt.Printf("✅ 查詢成功: %d\n", result)
		}

		// 測試表查詢
		fmt.Println("\n🧪 測試表查詢...")
		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public'").Scan(&count)
		if err != nil {
			fmt.Printf("❌ 表查詢失敗: %v\n", err)
		} else {
			fmt.Printf("✅ 表查詢成功: 找到 %d 個表\n", count)
		}

		// 測試 counseling_centers 表
		fmt.Println("\n🧪 測試 counseling_centers 表...")
		var centerCount int
		err = db.QueryRow("SELECT COUNT(*) FROM counseling_centers WHERE deleted_at IS NULL").Scan(&centerCount)
		if err != nil {
			fmt.Printf("❌ counseling_centers 查詢失敗: %v\n", err)
		} else {
			fmt.Printf("✅ counseling_centers 查詢成功: %d 筆記錄\n", centerCount)
		}

		db.Close()
		break
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// maskPassword 會將 DSN 字串中的密碼遮罩，避免敏感資訊外洩。
// 僅遮罩 password=xxx 或 user:password@ 這兩種常見格式。
func maskPassword(dsn string) string {
	// 嘗試遮罩 password=xxx 格式
	re := regexp.MustCompile(`(password=)([^&\s]+)`)
	masked := re.ReplaceAllString(dsn, "${1}***")

	// 嘗試遮罩 user:password@ 格式
	re2 := regexp.MustCompile(`([a-zA-Z0-9._%+-]+:)([^@]+)(@)`)
	masked = re2.ReplaceAllString(masked, "${1}***${3}")

	return masked
}
