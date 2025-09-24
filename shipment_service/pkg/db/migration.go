package db

import (
	"billing-system/shipment_service/internal/model"
	"log"

	"gorm.io/gorm"
)

// MigrateDB creates or updates the database schema
func MigrateDB(db *gorm.DB) error {
	log.Println("Running database migrations...")

	err := db.AutoMigrate(
		&model.Shipment{},
		&model.ShipmentItem{},
	)
	if err != nil {
		return err
	}

	log.Println("Database migrations completed successfully")
	return nil
}
