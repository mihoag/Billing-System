package db

import (
	"billing-system/billing_service/internal/model"
	"log"

	"gorm.io/gorm"
)

// MigrateDB creates or updates the database schema
func MigrateDB(db *gorm.DB) error {
	log.Println("Running database migrations...")

	err := db.AutoMigrate(
		&model.Item{},
		&model.Order{},
		&model.OrderItem{},
		&model.Payment{},
		&model.Invoice{},
		&model.InvoiceItem{},
	)
	if err != nil {
		return err
	}

	log.Println("Database migrations completed successfully")
	return nil
}
