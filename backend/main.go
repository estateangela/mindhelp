package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "mindhelp-backend/docs" // 導入 Swagger 文檔
	"mindhelp-backend/internal/config"
	"mindhelp-backend/internal/database"
	"mindhelp-backend/internal/routes"
)

// @title MindHelp Backend API
// @version 1.0
// @description MindHelp 心理健康支援應用程式後端 API
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host mindhelp.onrender.com
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description 請輸入 "Bearer " 加上 JWT token

func main() {
	// 載入配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 獲取端口 - Render 使用 PORT 環境變數
	port := os.Getenv("PORT")
	if port == "" {
		port = cfg.Server.Port
	}
	log.Printf("環境變數 PORT: %s", os.Getenv("PORT"))
	log.Printf("配置檔案端口: %s", cfg.Server.Port)
	log.Printf("最終使用端口: %s", port)

	// 設定路由 (不需要資料庫連接也能啟動基本路由)
	router := routes.SetupRoutes(cfg)

	// 創建 HTTP 伺服器
	srv := &http.Server{
		Addr:         "0.0.0.0:" + port, // 綁定到所有介面
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// 先啟動伺服器讓 Render 偵測到端口
	go func() {
		log.Printf("Starting server on 0.0.0.0:%s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// 在背景連接資料庫
	go func() {
		log.Println("Connecting to database in background...")
		if err := database.Connect(cfg); err != nil {
			log.Printf("Failed to connect to database: %v", err)
			log.Println(os.Getenv("DATABASE_URL"))
			// 不要讓資料庫連接失敗導致整個服務崩潰
			return
		}
		defer database.Close()

		// 執行資料庫遷移
		log.Println("Starting database migration...")
		if err := database.Migrate(); err != nil {
			log.Printf("Failed to migrate database: %v", err)
		} else {
			log.Println("Database migration completed successfully")
		}
	}()

	// 等待中斷信號
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// 優雅關閉伺服器
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited")
}
