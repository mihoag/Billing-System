package tests

import (
	"log"
	"regexp"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// MockDB holds the mock database connection and mock object
type MockDB struct {
	DB   *gorm.DB
	Mock sqlmock.Sqlmock
}

// NewMockDB creates a new mock database connection and mock object
func NewMockDB() (*MockDB, error) {
	// Create a new SQL mock database
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("Failed to create sqlmock: %v", err)
		return nil, err
	}

	// Configure GORM to use the mock database
	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_billing",
		DriverName:           "postgres",
		Conn:                 sqlDB,
		PreferSimpleProtocol: true,
	})

	// Create GORM DB with the mock database
	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatalf("Failed to open gorm DB: %v", err)
		return nil, err
	}

	return &MockDB{
		DB:   db,
		Mock: mock,
	}, nil
}

// Close closes the mock database connection
func (m *MockDB) Close() error {
	db, err := m.DB.DB()
	if err != nil {
		return err
	}
	return db.Close()
}

// ExpectationsWereMet checks if all expectations were met
func (m *MockDB) ExpectationsWereMet() error {
	return m.Mock.ExpectationsWereMet()
}

// Helper functions to create common mock column definitions
func ItemColumns() []string {
	return []string{"id", "created_at", "updated_at", "deleted_at", "name", "sku", "price"}
}

func OrderColumns() []string {
	return []string{"id", "created_at", "updated_at", "deleted_at", "customer_id", "total_amount", "status"}
}

func OrderItemColumns() []string {
	return []string{"id", "created_at", "updated_at", "deleted_at", "order_id", "quantity", "item_id"}
}

func PaymentColumns() []string {
	return []string{"id", "created_at", "updated_at", "deleted_at", "order_id", "method", "amount"}
}

func InvoiceColumns() []string {
	return []string{"id", "created_at", "updated_at", "deleted_at", "order_id", "shipment_id", "total_amount"}
}

func InvoiceItemColumns() []string {
	return []string{"id", "created_at", "updated_at", "deleted_at", "invoice_id", "quantity", "item_id"}
}

// Helper to convert Go time to SQL format
func AnyTime() sqlmock.Argument {
	return sqlmock.AnyArg()
}

// Helper for creating a regexp for SQL query matching
func QueryMatcher(query string) *regexp.Regexp {
	return regexp.MustCompile(query)
}
