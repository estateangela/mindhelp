package main

import (
	"log"
	"net/http"
	"time"

	"mindhelp-backend/internal/config"
	"mindhelp-backend/internal/database"
	"mindhelp-backend/internal/routes"
)

func main_test() {
	// 載入配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 連接到資料庫
	if err := database.Connect(cfg); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 設定路由
	r := routes.SetupRoutes(cfg)

	// 啟動伺服器
	address := ":" + cfg.Server.Port
	server := &http.Server{
		Addr:         address,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("🚀 Server starting on %s", address)
	log.Printf("📚 Swagger documentation: http://localhost%s/swagger/index.html", address)
	log.Printf("🔍 New API endpoints:")
	log.Printf("   - GET /api/v1/counselors")
	log.Printf("   - GET /api/v1/counseling-centers")
	log.Printf("   - GET /api/v1/recommended-doctors")

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to start server: %v", err)
	}
}
