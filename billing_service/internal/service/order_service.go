package service

import (
	"billing-system/billing_service/internal/dto"
	"billing-system/billing_service/internal/model"
	"billing-system/billing_service/internal/repository"
	"context"
	"fmt"
)

// OrderServiceImpl implements OrderService
type OrderServiceImpl struct {
	orderRepo repository.OrderRepository
	itemRepo  repository.ItemRepository
}

// NewOrderService creates a new OrderServiceImpl
func NewOrderService(
	orderRepo repository.OrderRepository,
	itemRepo repository.ItemRepository,
) OrderService {
	return &OrderServiceImpl{
		orderRepo: orderRepo,
		itemRepo:  itemRepo,
	}
}

// CreateOrder creates a new order with items and payments
func (s *OrderServiceImpl) CreateOrder(
	ctx context.Context,
	customerID string,
	itemRequests []dto.ItemRequest,
	paymentRequests []dto.PaymentRequest,
) (*model.Order, error) {
	// Calculate total amount from items
	var totalAmount float64
	orderItems := make([]model.OrderItem, 0, len(itemRequests))

	// Process items and calculate totals
	for _, req := range itemRequests {
		// Fetch item details from repository
		item, err := s.itemRepo.GetBySku(ctx, req.Sku)
		if err != nil {
			return nil, fmt.Errorf("item with SKU %s not found: %w", req.Sku, err)
		}

		totalAmount += float64(req.Quantity) * item.Price

		orderItems = append(orderItems, model.OrderItem{
			ItemID:   item.ID,
			Quantity: req.Quantity,
		})
	}

	// Calculate total payment amount
	var totalPayment float64
	payments := make([]model.Payment, 0, len(paymentRequests))

	for _, req := range paymentRequests {
		totalPayment += req.Amount
		payments = append(payments, model.Payment{
			Method: req.Method,
			Amount: req.Amount,
		})
	}

	// Validate that total payment equals total amount
	if totalPayment != totalAmount {
		return nil, fmt.Errorf("%w: payment total %f does not match order total %f",
			ErrInvalidAmount, totalPayment, totalAmount)
	}

	// Create order
	order := &model.Order{
		CustomerID:  customerID,
		TotalAmount: totalAmount,
		Status:      model.OrderPending,
		Items:       orderItems,
		Payments:    payments,
	}

	// Save the order to the database
	if err := s.orderRepo.Create(ctx, order); err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	// Return the created order
	return order, nil
}
