package main

import (
	"billing-system/bff/router"
)

func main() {

	// Load configuration from config.yaml
	// err := config.LoadConfig()
	// if err != nil {
	// 	log.Fatalf("Failed to load configuration: %v", err)
	// }
	// Start the router (this will block until the server shuts down)
	router.Start()
}
