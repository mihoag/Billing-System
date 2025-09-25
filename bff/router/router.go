package router

import (
	"log"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	billing "billing-system/bff/internal/billing"
	shipment "billing-system/bff/internal/shipment"
)

func Start() {
	// Create a default gin router
	router := gin.Default()

	// Initialize billing handler
	billingHandler := billing.NewHandler()
	shipmentHandler := shipment.NewHandler()

	// Set up billing API routes
	billingRoutes := router.Group("/api/v1")
	{
		// Order endpoints
		billingRoutes.POST("/orders", billingHandler.CreateOrder)
		billingRoutes.POST("/shipments", shipmentHandler.CreateShipment)
	}

	// Start HTTP server
	serverAddress := "127.0.0.1:8081" //config.Service.Server.Address

	if err := router.Run(serverAddress); err != nil {
		log.Fatal("Failed to start server", zap.Error(err))
	}
}
