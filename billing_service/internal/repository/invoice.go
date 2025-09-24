package repository

import (
	"billing-system/billing_service/internal/model"
	"context"
	"fmt"

	"gorm.io/gorm"
)

// InvoiceRepositoryImpl implements the InvoiceRepository interface
type InvoiceRepositoryImpl struct {
	db *gorm.DB
}

// NewItemRepository creates a new instance of ItemRepositoryImpl
func NewInvoiceRepository(db *gorm.DB) InvoiceRepository {
	return &InvoiceRepositoryImpl{
		db: db,
	}
}

// Create a new order in the database
func (r *InvoiceRepositoryImpl) Create(ctx context.Context, invoice *model.Invoice) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(invoice).Error; err != nil {
			return fmt.Errorf("failed to create invoice: %w", err)
		}
		return nil
	})
}
