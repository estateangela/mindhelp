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
			// 使用者管理路由
			users := protected.Group("/users")
			{
				userHandler := handlers.NewUserHandler(cfg)
				users.GET("/me", userHandler.GetProfile)
				users.PUT("/me", userHandler.UpdateProfile)
				users.DELETE("/me", userHandler.DeleteAccount)
				users.PUT("/me/password", userHandler.ChangePassword)
				users.GET("/me/stats", userHandler.GetStats)
			}

			// 聊天路由
			chat := protected.Group("/chat")
			{
				chatHandler := handlers.NewChatHandler(cfg)
				chat.POST("/send", chatHandler.SendMessage)
				chat.GET("/history", chatHandler.GetChatHistory)
				// TODO: 實現 session-based endpoints
				// chat.GET("/sessions", chatHandler.GetSessions)
				// chat.POST("/sessions", chatHandler.CreateSession)
				// chat.GET("/sessions/:sessionId/messages", chatHandler.GetSessionMessages)
				// chat.POST("/sessions/:sessionId/messages", chatHandler.SendSessionMessage)
			}

			// 位置路由 (需要認證的)
			locations := protected.Group("/locations")
			{
				locationHandler := handlers.NewLocationHandler()
				locations.POST("", locationHandler.CreateLocation)
				locations.PUT("/:id", locationHandler.UpdateLocation)
				locations.DELETE("/:id", locationHandler.DeleteLocation)
			}

			// 測驗路由 (需要認證的)
			quizzes := protected.Group("/quizzes")
			{
				quizHandler := handlers.NewQuizHandler()
				quizzes.POST("/:id/submit", quizHandler.SubmitQuiz)
			}

			// 使用者測驗歷史
			{
				quizHandler := handlers.NewQuizHandler()
				protected.GET("/users/me/quiz_history", quizHandler.GetQuizHistory)
			}

			// 收藏路由
			bookmarks := protected.Group("/users/me/bookmarks")
			{
				bookmarkHandler := handlers.NewBookmarkHandler()
				bookmarks.GET("/articles", bookmarkHandler.GetArticleBookmarks)
				bookmarks.GET("/resources", bookmarkHandler.GetLocationBookmarks)
			}

			// 通用收藏操作
			{
				bookmarkHandler := handlers.NewBookmarkHandler()
				protected.POST("/bookmarks", bookmarkHandler.BookmarkResource)
				protected.DELETE("/bookmarks", bookmarkHandler.UnbookmarkResource)
			}

			// 文章收藏 (舊版相容)
			{
				articleHandler := handlers.NewArticleHandler()
				protected.POST("/articles/:id/bookmark", articleHandler.BookmarkArticle)
				protected.DELETE("/articles/:id/bookmark", articleHandler.UnbookmarkArticle)
			}
		}

		// 公開路由
		{
			// 位置相關公開路由
			locationHandler := handlers.NewLocationHandler()
			api.GET("/locations/search", locationHandler.SearchLocations)
			api.GET("/locations/:id", locationHandler.GetLocation)

			// 文章相關公開路由
			articleHandler := handlers.NewArticleHandler()
			api.GET("/articles", articleHandler.GetArticles)
			api.GET("/articles/:id", articleHandler.GetArticle)

			// 測驗相關公開路由
			quizHandler := handlers.NewQuizHandler()
			api.GET("/quizzes", quizHandler.GetQuizzes)
			api.GET("/quizzes/:id", quizHandler.GetQuiz)

			// 應用配置
			configHandler := handlers.NewConfigHandler()
			api.GET("/config", configHandler.GetConfig)
		}
	}

	return r
}
