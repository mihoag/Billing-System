package repository

import (
	"billing-system/billing_service/internal/model"
	"context"
	"errors"
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

		// If there are order items, create them
		if len(order.Items) > 0 {
			for i := range order.Items {
				order.Items[i].OrderID = order.ID // Set the order ID for each item
			}
			if err := tx.Create(&order.Items).Error; err != nil {
				return fmt.Errorf("failed to create order items: %w", err)
			}
		}

		// If there are payments, create them
		if len(order.Payments) > 0 {
			for i := range order.Payments {
				order.Payments[i].OrderID = order.ID // Set the order ID for each payment
			}
			if err := tx.Create(&order.Payments).Error; err != nil {
				return fmt.Errorf("failed to create order payments: %w", err)
			}
		}

		return nil
	})
}

// GetByID retrieves an order by its ID from the database
func (r *OrderRepositoryImpl) GetByID(ctx context.Context, id int64) (*model.Order, error) {
	var order model.Order

	// Query the order with preloaded relationships
	result := r.db.WithContext(ctx).
		Preload("Items").
		Preload("Payments").
		Preload("Invoices").
		First(&order, id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("order with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to get order: %w", result.Error)
	}

	return &order, nil
}
