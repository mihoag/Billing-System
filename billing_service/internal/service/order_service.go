package service

import (
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
) *OrderServiceImpl {
	return &OrderServiceImpl{
		orderRepo: orderRepo,
		itemRepo:  itemRepo,
	}
}

// CreateOrder creates a new order with items and payments
func (s *OrderServiceImpl) CreateOrder(
	ctx context.Context,
	customerID string,
	itemRequests []ItemRequest,
	paymentRequests []PaymentRequest,
) (*model.Order, error) {
	// Calculate total amount from items
	var totalAmount float64
	orderItems := make([]model.OrderItem, 0, len(itemRequests))

	// Process items and calculate totals
	for _, req := range itemRequests {
		item, err := s.itemRepo.GetByID(ctx, req.ItemID)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrItemNotFound, err)
		}

		// Use provided price if available, otherwise use item's price
		unitPrice := item.Price
		if req.Price > 0 {
			unitPrice = req.Price
		}

		total := float64(req.Quantity) * unitPrice
		totalAmount += total

		orderItems = append(orderItems, model.OrderItem{
			ItemId:   item.ID,
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
		Status:      model.Pending,
		Items:       orderItems,
		Payments:    payments,
	}

	// Save the order to the database
	if err := s.orderRepo.Create(ctx, order); err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	// If we got this far, the order was created successfully
	return order, nil
}

// GetOrder retrieves an order by ID with all associations
func (s *OrderServiceImpl) GetOrder(ctx context.Context, id int64) (*model.Order, error) {
	order, err := s.orderRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrOrderNotFound, err)
	}
	return order, nil
}
