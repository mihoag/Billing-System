package tests

import (
	"billing-system/billing_service/internal/dto"
	"billing-system/billing_service/internal/model"
	"billing-system/billing_service/internal/service"
	"billing-system/billing_service/internal/service/tests/mocks"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestOrderService_CreateOrder(t *testing.T) {
	// Define test time for consistent timestamps
	testTime := time.Now()

	// Define test cases
	testCases := []struct {
		name            string
		customerID      string
		itemRequests    []dto.ItemRequest
		paymentRequests []dto.PaymentRequest
		mockSetup       func(*mocks.MockOrderRepository, *mocks.MockItemRepository)
		expectedError   error
		checkOrder      func(*testing.T, *model.Order)
	}{
		{
			name:       "Success - Create order with valid items and payments",
			customerID: "customer-123",
			itemRequests: []dto.ItemRequest{
				{Sku: "SKU001", Quantity: 2},
				{Sku: "SKU002", Quantity: 1},
			},
			paymentRequests: []dto.PaymentRequest{
				{Method: "COD", Amount: 300},
			},
			mockSetup: func(orderRepo *mocks.MockOrderRepository, itemRepo *mocks.MockItemRepository) {
				// Mock item repository responses
				itemRepo.On("GetBySku", mock.Anything, "SKU001").Return(&model.Item{
					Base:  model.Base{ID: 1},
					Name:  "Item 1",
					Sku:   "SKU001",
					Price: 100,
				}, nil)

				itemRepo.On("GetBySku", mock.Anything, "SKU002").Return(&model.Item{
					Base:  model.Base{ID: 2},
					Name:  "Item 2",
					Sku:   "SKU002",
					Price: 100,
				}, nil)

				// Mock order repository create
				orderRepo.On("Create", mock.Anything, mock.MatchedBy(func(order *model.Order) bool {
					return order.CustomerID == "customer-123" &&
						order.TotalAmount == 300 &&
						order.Status == model.OrderPending &&
						len(order.Items) == 2 &&
						len(order.Payments) == 1
				})).Return(nil).Run(func(args mock.Arguments) {
					// Set ID on the order when created
					order := args.Get(1).(*model.Order)
					order.ID = 1
					order.CreatedAt = testTime
					order.UpdatedAt = testTime

					// Set Order ID on items and payments
					for i := range order.Items {
						order.Items[i].OrderID = 1
						order.Items[i].ID = int64(i + 1)
					}
					for i := range order.Payments {
						order.Payments[i].OrderID = 1
						order.Payments[i].ID = int64(i + 1)
					}
				})
			},
			expectedError: nil,
			checkOrder: func(t *testing.T, order *model.Order) {
				assert.Equal(t, int64(1), order.ID)
				assert.Equal(t, "customer-123", order.CustomerID)
				assert.Equal(t, 300.0, order.TotalAmount)
				assert.Equal(t, model.OrderPending, order.Status)

				// Check items
				assert.Len(t, order.Items, 2)
				assert.Equal(t, int64(1), order.Items[0].ItemID)
				assert.Equal(t, 2, order.Items[0].Quantity)
				assert.Equal(t, int64(2), order.Items[1].ItemID)
				assert.Equal(t, 1, order.Items[1].Quantity)

				// Check payments
				assert.Len(t, order.Payments, 1)
				assert.Equal(t, model.COD, order.Payments[0].Method)
				assert.Equal(t, 300.0, order.Payments[0].Amount)
			},
		},
		{
			name:       "Error - Item not found",
			customerID: "customer-123",
			itemRequests: []dto.ItemRequest{
				{Sku: "INVALID-SKU", Quantity: 1},
			},
			paymentRequests: []dto.PaymentRequest{
				{Method: "COD", Amount: 99.99},
			},
			mockSetup: func(orderRepo *mocks.MockOrderRepository, itemRepo *mocks.MockItemRepository) {
				// Mock item repository to return error for invalid SKU
				itemRepo.On("GetBySku", mock.Anything, "INVALID-SKU").Return(nil, service.ErrItemNotFound)
			},
			expectedError: service.ErrItemNotFound,
			checkOrder:    nil,
		},
		{
			name:       "Error - Payment amount does not match order total",
			customerID: "customer-123",
			itemRequests: []dto.ItemRequest{
				{Sku: "SKU001", Quantity: 2},
			},
			paymentRequests: []dto.PaymentRequest{
				{Method: "COD", Amount: 150.00}, // Incorrect amount
			},
			mockSetup: func(orderRepo *mocks.MockOrderRepository, itemRepo *mocks.MockItemRepository) {
				// Mock item repository response
				itemRepo.On("GetBySku", mock.Anything, "SKU001").Return(&model.Item{
					Base:  model.Base{ID: 1},
					Name:  "Item 1",
					Sku:   "SKU001",
					Price: 99.99, // Total should be 199.98
				}, nil)
			},
			expectedError: service.ErrInvalidAmount,
			checkOrder:    nil,
		},
		{
			name:       "Error - Database error during order creation",
			customerID: "customer-123",
			itemRequests: []dto.ItemRequest{
				{Sku: "SKU001", Quantity: 1},
			},
			paymentRequests: []dto.PaymentRequest{
				{Method: "COD", Amount: 99.99},
			},
			mockSetup: func(orderRepo *mocks.MockOrderRepository, itemRepo *mocks.MockItemRepository) {
				// Mock item repository response
				itemRepo.On("GetBySku", mock.Anything, "SKU001").Return(&model.Item{
					Base:  model.Base{ID: 1},
					Name:  "Item 1",
					Sku:   "SKU001",
					Price: 99.99,
				}, nil)

				// Mock order repository to return error
				orderRepo.On("Create", mock.Anything, mock.AnythingOfType("*model.Order")).Return(errors.New("database error"))
			},
			expectedError: errors.New("database error"),
			checkOrder:    nil,
		},
		{
			name:         "Success - Zero items (edge case)",
			customerID:   "customer-123",
			itemRequests: []dto.ItemRequest{},
			paymentRequests: []dto.PaymentRequest{
				{Method: "COD", Amount: 0.0},
			},
			mockSetup: func(orderRepo *mocks.MockOrderRepository, itemRepo *mocks.MockItemRepository) {
				// Mock order repository create for empty order
				orderRepo.On("Create", mock.Anything, mock.MatchedBy(func(order *model.Order) bool {
					return order.CustomerID == "customer-123" &&
						order.TotalAmount == 0.0 &&
						len(order.Items) == 0 &&
						len(order.Payments) == 1
				})).Return(nil).Run(func(args mock.Arguments) {
					// Set ID on the order when created
					order := args.Get(1).(*model.Order)
					order.ID = 1
				})
			},
			expectedError: nil,
			checkOrder: func(t *testing.T, order *model.Order) {
				assert.Equal(t, int64(1), order.ID)
				assert.Equal(t, 0.0, order.TotalAmount)
				assert.Empty(t, order.Items)
				assert.Len(t, order.Payments, 1)
				assert.Equal(t, 0.0, order.Payments[0].Amount)
			},
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create mocks
			mockOrderRepo := new(mocks.MockOrderRepository)
			mockItemRepo := new(mocks.MockItemRepository)

			// Set up mocks
			tc.mockSetup(mockOrderRepo, mockItemRepo)

			// Create service with mocks
			orderService := service.NewOrderService(mockOrderRepo, mockItemRepo)

			// Call the method being tested
			order, err := orderService.CreateOrder(context.Background(), tc.customerID, tc.itemRequests, tc.paymentRequests)

			// Check errors
			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tc.expectedError.Error())
				assert.Nil(t, order)
			} else {
				assert.NoError(t, err)
				require.NotNil(t, order)
				if tc.checkOrder != nil {
					tc.checkOrder(t, order)
				}
			}

			// Verify that all expected mock calls were made
			mockOrderRepo.AssertExpectations(t)
			mockItemRepo.AssertExpectations(t)
		})
	}
}
