package router

import (
	"log"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	billing "billing-system/bff/internal/billing"
)

func Start() {
	// Create a default gin router
	router := gin.Default()

	// Initialize billing handler
	billingHandler := billing.NewHandler()

	// Set up billing API routes
	billingRoutes := router.Group("/api/v1/billing")
	{
		// Order endpoints
		billingRoutes.POST("/orders", billingHandler.CreateOrder)
	}

	// Start HTTP server
	serverAddress := "127.0.0.1:8081" //config.Service.Server.Address

	if err := router.Run(serverAddress); err != nil {
		log.Fatal("Failed to start server", zap.Error(err))
	}
}
