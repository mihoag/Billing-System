package main

import (
	"go.uber.org/zap"

	"billing-system/bff/config"
	"billing-system/bff/router"
)

func main() {
	// Initialize logger
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// Load configuration from config.yaml
	err := config.LoadConfig()
	if err != nil {
		logger.Fatal("Failed to load config", zap.Error(err))
	}

	// Use the configured logger
	config.Service.Logger = logger
	logger.Info("Configuration loaded successfully")

	// Start the router (this will block until the server shuts down)
	router.Start()
}
