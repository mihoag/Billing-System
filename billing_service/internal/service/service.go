package service

import (
	"billing-system/billing_service/internal/model"
	"context"
	"errors"
)

var (
	ErrItemNotFound        = errors.New("item not found")
	ErrOrderNotFound       = errors.New("order not found")
	ErrInvoiceNotFound     = errors.New("invoice not found")
	ErrPaymentNotFound     = errors.New("payment not found")
	ErrInvalidQuantity     = errors.New("invalid quantity")
	ErrInvalidAmount       = errors.New("invalid amount")
	ErrInsufficientPayment = errors.New("insufficient payment")
	ErrDatabaseError       = errors.New("database error")
)

// OrderService defines the interface for order-related business logic
type OrderService interface {
	CreateOrder(ctx context.Context, customerID string, items []ItemRequest, payments []PaymentRequest) (*model.Order, error)
	GetOrder(ctx context.Context, id int64) (*model.Order, error)
}

// InvoiceService defines the interface for invoice-related business logic
type InvoiceService interface {
	CreateInvoice(ctx context.Context, orderID int64, shipmentID int64, items []ItemRequest) (*model.Invoice, error)
	GetInvoice(ctx context.Context, id int64) (*model.Invoice, error)
}

// ItemRequest represents a request to include an item in an order or invoice
type ItemRequest struct {
	ItemID   int64   `json:"item_id"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price,omitempty"` // Optional for custom pricing
}

// PaymentRequest represents a request to add a payment to an order
type PaymentRequest struct {
	Method model.PaymentMethod `json:"method"`
	Amount float64             `json:"amount"`
}
