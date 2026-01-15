package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/lanhyde/ogenkidesuka-server/internal/config"
	"github.com/lanhyde/ogenkidesuka-server/internal/database"
	"github.com/lanhyde/ogenkidesuka-server/internal/handlers"
	"github.com/lanhyde/ogenkidesuka-server/internal/middleware"
	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	if err := database.Connect(&cfg.Database); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	defer database.Close()

	// Set up Gin router
	if cfg.Server.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	// Apply middleware
	router.Use(middleware.CORS(cfg.CORS.AllowedOrigins))

	// API routes
	api := router.Group("/api/v1")
	{
		// Health check
		api.GET("/health", handlers.HealthCheck)

		// Check-in routes
		checkins := api.Group("/checkins")
		{
			checkins.POST("/:userId", handlers.CreateCheckIn)
			checkins.GET("/:userId/today", handlers.GetTodayCheckIn)
			checkins.GET("/:userId/history", handlers.GetCheckInHistory)
		}
	}

	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	fmt.Printf("\n🚀 Server starting on http://localhost%s\n", addr)
	fmt.Printf("📊 Environment: %s\n", cfg.Server.Env)
	fmt.Printf("💾 Database: %s@%s:%s/%s\n\n",
		cfg.Database.User,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.DBName,
	)
	if err := router.Run(addr); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
