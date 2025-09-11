package routes

import (
	"mindhelp-backend/internal/config"
	"mindhelp-backend/internal/handlers"
	"mindhelp-backend/internal/middleware"
	"os"

	"github.com/gin-gonic/gin"
	ginCors "github.com/rs/cors/wrapper/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRoutes 設定路由
func SetupRoutes(cfg *config.Config) *gin.Engine {
	// 設定 Gin 模式
	gin.SetMode(cfg.Server.GinMode)

	// 創建路由引擎
	r := gin.New()

	// 使用中間件
	logFile := os.Getenv("LOG_FILE")
	r.Use(middleware.StructuredLogger(logFile))
	r.Use(middleware.MetricsMiddleware())
	r.Use(gin.Recovery())

	// CORS 中間件
	corsConfig := ginCors.New(ginCors.Options{
		AllowedOrigins:   cfg.CORS.AllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	})
	r.Use(corsConfig)

	// Swagger 文檔
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 監控處理器
	monitoringHandler := handlers.NewMonitoringHandler()

	// 健康檢查和監控端點
	r.GET("/health", monitoringHandler.HealthCheck)
	r.GET("/health/detailed", monitoringHandler.DetailedHealthCheck)
	r.GET("/health/ready", monitoringHandler.ReadinessCheck)
	r.GET("/health/live", monitoringHandler.LivenessCheck)
	r.GET("/metrics", monitoringHandler.Metrics)

	// API 路由組
	api := r.Group("/api/v1")
	{
		// 認證路由
		auth := api.Group("/auth")
		{
			authHandler := handlers.NewAuthHandler(cfg)
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
		}

		// 需要認證的路由
		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware(cfg))
		{
			// 聊天路由
			chat := protected.Group("/chat")
			{
				chatHandler := handlers.NewChatHandler(cfg)
				chat.POST("/send", chatHandler.SendMessage)
				chat.GET("/history", chatHandler.GetChatHistory)
			}

			// 位置路由
			locations := protected.Group("/locations")
			{
				locationHandler := handlers.NewLocationHandler()
				locations.POST("", locationHandler.CreateLocation)
				locations.PUT("/:id", locationHandler.UpdateLocation)
				locations.DELETE("/:id", locationHandler.DeleteLocation)
			}
		}

		// 公開路由
		{
			locationHandler := handlers.NewLocationHandler()
			api.GET("/locations/search", locationHandler.SearchLocations)
			api.GET("/locations/:id", locationHandler.GetLocation)
		}
	}

	return r
}
