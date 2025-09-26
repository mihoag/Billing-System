package db

import (
	"fmt"

	"billing-system/billing_service/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Connect connects to the database and returns a gorm.DB instance
func NewDatabase(config *config.Config) (*gorm.DB, error) {
	dbConfig := config.Database
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s ",
		dbConfig.Host, dbConfig.Port, dbConfig.Username, dbConfig.Password, dbConfig.Database,
	)

	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}
