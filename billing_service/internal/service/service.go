package service

import (
	"billing-system/billing_service/internal/model"
	pb "billing-system/billing_service/proto"
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
}

type InvoiceService interface {
	CreateInvoice(ctx context.Context, createInvoiceRequest *pb.CreateInvoiceRequest) (*model.Invoice, error)
}

// ItemRequest represents a request to include an item in an order or invoice
type ItemRequest struct {
	Sku      string `json:"skus"`
	Quantity int    `json:"quantity"`
}

// PaymentRequest represents a request to add a payment to an order
type PaymentRequest struct {
	Method model.PaymentMethod `json:"method"`
	Amount float64             `json:"amount"`
}
