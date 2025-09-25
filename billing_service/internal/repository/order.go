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
