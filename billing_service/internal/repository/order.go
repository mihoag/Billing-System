package repository

import (
	"billing-system/billing_service/internal/model"
	"context"

	"gorm.io/gorm"
)

// OrderRepositoryImpl implements the OrderRepository interface
type OrderRepositoryImpl struct {
	db *gorm.DB
}

// NewOrderRepository creates a new instance of OrderRepositoryImpl
func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &OrderRepositoryImpl{
		db: db,
	}
}

// Create a new order in the database along with its associated items and payments.
// It uses a transaction to ensure all data is saved atomically.
// Returns an error if the creation fails.
func (r *OrderRepositoryImpl) Create(ctx context.Context, order *model.Order) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(order).Error; err != nil {
			return err
		}
		return nil
	})
}

// GetByID retrieves an order by its ID along with related items and payments.
// Returns the order and nil if found, nil and error otherwise.
func (r *OrderRepositoryImpl) GetByID(ctx context.Context, id int64) (*model.Order, error) {
	var order model.Order

	result := r.db.WithContext(ctx).
		Preload("Items.Item"). // Preload items and their details
		Preload("Payments").   // Preload payment information
		First(&order, id)      // Find by primary key

	if result.Error != nil {
		return nil, result.Error // Return the error (could be gorm.ErrRecordNotFound)
	}

	return &order, nil
}
