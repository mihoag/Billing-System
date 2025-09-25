package repository

import (
	"billing-system/billing_service/internal/model"
	"context"

	"gorm.io/gorm"
)

// InvoiceRepositoryImpl implements the InvoiceRepository interface
type InvoiceRepositoryImpl struct {
	db *gorm.DB
}

// NewInvoiceRepository creates a new instance of InvoiceRepositoryImpl
func NewInvoiceRepository(db *gorm.DB) InvoiceRepository {
	return &InvoiceRepositoryImpl{
		db: db,
	}
}

// Create a new invoice in the database
func (r *InvoiceRepositoryImpl) Create(ctx context.Context, invoice *model.Invoice) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(invoice).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *InvoiceRepositoryImpl) GetByOrderID(ctx context.Context, orderID int64) ([]model.Invoice, error) {
	var invoices []model.Invoice

	result := r.db.WithContext(ctx).
		Where("order_id = ?", orderID).
		Preload("Items").
		Find(&invoices)

	if result.Error != nil {
		return nil, result.Error
	}

	return invoices, nil
}
