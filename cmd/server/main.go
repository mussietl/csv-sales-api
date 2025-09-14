package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/mussietl/csv-sales-api/internal/handlers"
	"github.com/mussietl/csv-sales-api/internal/services"
	"github.com/mussietl/csv-sales-api/pkg/utils"
	"github.com/sirupsen/logrus"
)

func main() {
	// Initialize logger
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.InfoLevel)

	// Create uploads directory if it doesn't exist
	uploadsDir := "public/uploads"
	if err := os.MkdirAll(uploadsDir, 0755); err != nil {
		logger.Fatalf("Failed to create uploads directory: %v", err)
	}

	// Initialize services
	fileService := services.NewFileService(uploadsDir, logger)
	csvService := services.NewCSVService(logger)

	// Initialize handlers
	uploadHandler := handlers.NewUploadHandler(fileService, csvService, logger)

	// Setup router
	router := gin.Default()

	// Add CORS middleware
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Routes
	api := router.Group("/api/v1")
	{
		api.POST("/upload", uploadHandler.UploadCSV)
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})
	}

	// Serve static files
	router.Static("/public", "./public")

	// Get port from environment or use default
	port := utils.GetEnv("PORT", "8080")

	logger.Infof("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
