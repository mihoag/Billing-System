package service

import (
	"billing-system/billing_service/internal/dto"
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
	CreateOrder(ctx context.Context, customerID string, items []dto.ItemRequest, payments []dto.PaymentRequest) (*model.Order, error)
}

type InvoiceService interface {
	CreateInvoice(ctx context.Context, shipmentId int64, orderId int64, itemRequest []dto.InvoiceItemRequest) (*model.Invoice, error)
}
