package repository

import (
	"billing-system/billing_service/internal/model"
	"context"
)

type ItemRepository interface {
	GetByID(ctx context.Context, id int64) (*model.Item, error)
}

// OrderRepository defines the interface for order operations
type OrderRepository interface {
	Create(ctx context.Context, order *model.Order) error
	GetByID(ctx context.Context, id int64) (*model.Order, error)
}

// OrderItemRepository defines the interface for order item operations
type OrderItemRepository interface {
	Create(ctx context.Context, orderItem *model.OrderItem) error
}

// PaymentRepository defines the interface for payment operations
type PaymentRepository interface {
	Create(ctx context.Context, payment *model.Payment) error
}

// InvoiceRepository defines the interface for invoice operations
type InvoiceRepository interface {
	Create(ctx context.Context, invoice *model.Invoice) error
	GetByID(ctx context.Context, id int64) (*model.Invoice, error)
	GetByOrderID(ctx context.Context, orderID int64) ([]*model.Invoice, error)
}

// InvoiceItemRepository defines the interface for invoice item operations
type InvoiceItemRepository interface {
	Create(ctx context.Context, invoiceItem *model.InvoiceItem) error
	GetByID(ctx context.Context, id int64) (*model.InvoiceItem, error)
	GetByInvoiceID(ctx context.Context, invoiceID int64) ([]*model.InvoiceItem, error)
}
