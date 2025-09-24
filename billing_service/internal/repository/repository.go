package repository

import (
	"billing-system/billing_service/internal/model"
	"context"
)

type ItemRepository interface {
	GetBySku(ctx context.Context, sku string) (*model.Item, error)
}

// OrderRepository defines the interface for order operations
type OrderRepository interface {
	Create(ctx context.Context, order *model.Order) error
}

// InvoiceRepository defines the interface for invoice operations
type InvoiceRepository interface {
	Create(ctx context.Context, invoice *model.Invoice) error
}
