package repository

import (
	"billing-system/billing_service/internal/model"
	"context"
	"fmt"

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

// Create creates a new order in the database
func (r *OrderRepositoryImpl) Create(ctx context.Context, order *model.Order) error {
	// Use a transaction to ensure atomicity when creating the order and related items
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Create the order first
		if err := tx.Create(order).Error; err != nil {
			return fmt.Errorf("failed to create order: %w", err)
		}

		return nil
	})
}
