package main

import (
	"log"

	"billing-system/shipment_service/pkg/config"
	"billing-system/shipment_service/pkg/db"
)

func main() {
	// load configuration from config.yaml
	// err := config.LoadConfig()
	// if err != nil {
	// 	log.Fatalf("Failed to get config: %v", err)
	// }

	gormDB, err := db.NewDatabase(&config.Service)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Run migrations
	if err := db.MigrateDB(gormDB); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

}
