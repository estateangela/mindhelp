package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	fmt.Println("🔍 Render 環境變數診斷工具")
	fmt.Println(strings.Repeat("=", 50))

	// 檢查 DATABASE_URL
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL != "" {
		fmt.Printf("✅ DATABASE_URL 已設定\n")
		fmt.Printf("📋 值: %s\n", maskDSN(databaseURL))
		fmt.Printf("🎯 狀態: 建議使用此配置\n")
	} else {
		fmt.Printf("❌ DATABASE_URL 未設定\n")
		fmt.Printf("💡 建議: 設定完整的 Supabase 連接字串\n")
	}

	fmt.Println()

	// 檢查個別資料庫參數
	dbParams := map[string]string{
		"DB_HOST":     os.Getenv("DB_HOST"),
		"DB_PORT":     os.Getenv("DB_PORT"),
		"DB_USER":     os.Getenv("DB_USER"),
		"DB_PASSWORD": os.Getenv("DB_PASSWORD"),
		"DB_NAME":     os.Getenv("DB_NAME"),
		"DB_SSL_MODE": os.Getenv("DB_SSL_MODE"),
	}

	fmt.Println("📊 個別資料庫參數檢查:")
	allParamsSet := true
	for key, value := range dbParams {
		if value != "" {
			if key == "DB_PASSWORD" {
				fmt.Printf("   ✅ %s: ***\n", key)
			} else {
				fmt.Printf("   ✅ %s: %s\n", key, value)
			}
		} else {
			fmt.Printf("   ❌ %s: 未設定\n", key)
			allParamsSet = false
		}
	}

	fmt.Println()

	if databaseURL != "" {
		fmt.Printf("🎯 推薦配置: 使用 DATABASE_URL\n")
		fmt.Printf("🔧 Render 環境變數設定:\n")
		fmt.Printf("   DATABASE_URL=%s\n", databaseURL)
	} else if allParamsSet {
		fmt.Printf("🎯 推薦配置: 使用個別參數\n")
		fmt.Printf("🔧 Render 環境變數設定:\n")
		fmt.Printf("   DB_HOST=%s\n", dbParams["DB_HOST"])
		fmt.Printf("   DB_PORT=%s\n", dbParams["DB_PORT"])
		fmt.Printf("   DB_USER=%s\n", dbParams["DB_USER"])
		fmt.Printf("   DB_PASSWORD=%s\n", dbParams["DB_PASSWORD"])
		fmt.Printf("   DB_NAME=%s\n", dbParams["DB_NAME"])
		fmt.Printf("   DB_SSL_MODE=%s\n", dbParams["DB_SSL_MODE"])
	} else {
		fmt.Printf("❌ 配置不完整\n")
		fmt.Printf("💡 請設定以下其中一種配置:\n")
		fmt.Printf("\n選項 A (推薦):\n")
		fmt.Printf("   DATABASE_URL=postgresql://postgres.haunuvdhisdygfradaya:MIND_HELP_2025@aws-1-ap-southeast-1.pooler.supabase.com:6543/postgres?sslmode=require&connect_timeout=30&pool_mode=transaction\n")
		fmt.Printf("\n選項 B:\n")
		fmt.Printf("   DB_HOST=aws-1-ap-southeast-1.pooler.supabase.com\n")
		fmt.Printf("   DB_PORT=6543\n")
		fmt.Printf("   DB_USER=postgres.haunuvdhisdygfradaya\n")
		fmt.Printf("   DB_PASSWORD=MIND_HELP_2025\n")
		fmt.Printf("   DB_NAME=postgres\n")
		fmt.Printf("   DB_SSL_MODE=require\n")
	}

	fmt.Println()
	fmt.Println("📝 設定步驟:")
	fmt.Println("1. 登入 Render Dashboard")
	fmt.Println("2. 進入你的服務設定")
	fmt.Println("3. 點擊 Environment 標籤")
	fmt.Println("4. 設定上述環境變數")
	fmt.Println("5. 點擊 Save changes")
	fmt.Println("6. 重新部署服務")
}

func maskDSN(dsn string) string {
	parts := strings.Split(dsn, " ")
	for i, part := range parts {
		if strings.HasPrefix(part, "password=") {
			parts[i] = "password=***"
			break
		}
	}
	return strings.Join(parts, " ")
}
