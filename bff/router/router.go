package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"billing-system/bff/config"
	billing "billing-system/bff/internal/billing"
)

func Start() {
	// Create a default gin router
	router := gin.Default()

	// Recovery middleware
	router.Use(gin.Recovery())

	// Add health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// Initialize billing handler
	billingHandler := billing.NewHandler()

	// Set up billing API routes
	billingRoutes := router.Group("/api/v1/billing")
	{
		// Order endpoints
		billingRoutes.POST("/orders", billingHandler.CreateOrder)
		billingRoutes.GET("/orders/:id", billingHandler.GetOrder)
	}

	// Start HTTP server
	serverAddress := config.Service.Server.Address
	config.Service.Logger.Info("Starting HTTP server", zap.String("address", serverAddress))
	if err := router.Run(serverAddress); err != nil {
		config.Service.Logger.Fatal("Failed to start HTTP server", zap.Error(err))
	}
}
