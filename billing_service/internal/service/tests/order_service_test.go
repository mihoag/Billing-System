package tests

import (
	"billing-system/billing_service/internal/dto"
	"billing-system/billing_service/internal/model"
	"billing-system/billing_service/internal/service"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock implementations of repositories
type MockOrderRepository struct {
	mock.Mock
}

func (m *MockOrderRepository) Create(ctx context.Context, order *model.Order) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

func (m *MockOrderRepository) GetByID(ctx context.Context, id int64) (*model.Order, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Order), args.Error(1)
}

func (m *MockOrderRepository) Update(ctx context.Context, order *model.Order) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

func (m *MockOrderRepository) GetAll(ctx context.Context) ([]model.Order, error) {
	args := m.Called(ctx)
	return args.Get(0).([]model.Order), args.Error(1)
}

func (m *MockOrderRepository) GetByCustomerID(ctx context.Context, customerID string) ([]model.Order, error) {
	args := m.Called(ctx, customerID)
	return args.Get(0).([]model.Order), args.Error(1)
}

type MockItemRepository struct {
	mock.Mock
}

func (m *MockItemRepository) GetBySku(ctx context.Context, sku string) (*model.Item, error) {
	args := m.Called(ctx, sku)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Item), args.Error(1)
}

func (m *MockItemRepository) Create(ctx context.Context, item *model.Item) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *MockItemRepository) GetByID(ctx context.Context, id int64) (*model.Item, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Item), args.Error(1)
}

func TestCreateOrder(t *testing.T) {
	// Define test cases using table-driven approach
	testCases := []struct {
		name          string
		customerID    string
		items         []dto.ItemRequest
		payments      []dto.PaymentRequest
		mockSetup     func(*MockOrderRepository, *MockItemRepository)
		expectedError error
		expectedOrder *model.Order
		checkOrder    func(*testing.T, *model.Order)
	}{
		{
			name:       "Success - Create order with valid inputs",
			customerID: "customer123",
			items: []dto.ItemRequest{
				{Sku: "SKU001", Quantity: 2},
				{Sku: "SKU002", Quantity: 1},
			},
			payments: []dto.PaymentRequest{
				{Method: model.COD, Amount: 150.0},
			},
			mockSetup: func(orderRepo *MockOrderRepository, itemRepo *MockItemRepository) {
				// Mock item repository responses
				itemRepo.On("GetBySku", mock.Anything, "SKU001").Return(&model.Item{
					Base:  model.Base{ID: 1},
					Sku:   "SKU001",
					Name:  "Item 1",
					Price: 50.0,
				}, nil)

				itemRepo.On("GetBySku", mock.Anything, "SKU002").Return(&model.Item{
					Base:  model.Base{ID: 2},
					Sku:   "SKU002",
					Name:  "Item 2",
					Price: 50.0,
				}, nil)

				// Mock successful order creation
				orderRepo.On("Create", mock.Anything, mock.AnythingOfType("*model.Order")).Run(func(args mock.Arguments) {
					// Set ID when order is created to simulate database behavior
					order := args.Get(1).(*model.Order)
					order.ID = 1
					order.CreatedAt = time.Now()
					order.UpdatedAt = time.Now()

					// Set order IDs for nested items
					for i := range order.Items {
						order.Items[i].OrderID = order.ID
					}

					// Set order IDs for payments
					for i := range order.Payments {
						order.Payments[i].OrderID = order.ID
					}
				}).Return(nil)
			},
			expectedError: nil,
			checkOrder: func(t *testing.T, order *model.Order) {
				assert.NotNil(t, order)
				assert.Equal(t, int64(1), order.ID)
				assert.Equal(t, "customer123", order.CustomerID)
				assert.Equal(t, 150.0, order.TotalAmount)
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
				assert.Equal(t, 150.0, order.Payments[0].Amount)
			},
		},
		{
			name:       "Failure - Item not found",
			customerID: "customer123",
			items: []dto.ItemRequest{
				{Sku: "SKU001", Quantity: 2},
				{Sku: "NONEXISTENT", Quantity: 1},
			},
			payments: []dto.PaymentRequest{
				{Method: model.COD, Amount: 150.0},
			},
			mockSetup: func(orderRepo *MockOrderRepository, itemRepo *MockItemRepository) {
				// Mock item repository responses
				itemRepo.On("GetBySku", mock.Anything, "SKU001").Return(&model.Item{
					Base:  model.Base{ID: 1},
					Sku:   "SKU001",
					Name:  "Item 1",
					Price: 50.0,
				}, nil)

				// Mock item not found error
				itemRepo.On("GetBySku", mock.Anything, "NONEXISTENT").Return(nil, service.ErrItemNotFound)
			},
			expectedError: service.ErrItemNotFound,
			checkOrder:    nil,
		},
		{
			name:       "Failure - Payment amount mismatch",
			customerID: "customer123",
			items: []dto.ItemRequest{
				{Sku: "SKU001", Quantity: 2},
				{Sku: "SKU002", Quantity: 1},
			},
			payments: []dto.PaymentRequest{
				{Method: model.COD, Amount: 100.0}, // Incorrect amount (should be 150.0)
			},
			mockSetup: func(orderRepo *MockOrderRepository, itemRepo *MockItemRepository) {
				// Mock item repository responses
				itemRepo.On("GetBySku", mock.Anything, "SKU001").Return(&model.Item{
					Base:  model.Base{ID: 1},
					Sku:   "SKU001",
					Name:  "Item 1",
					Price: 50.0,
				}, nil)

				itemRepo.On("GetBySku", mock.Anything, "SKU002").Return(&model.Item{
					Base:  model.Base{ID: 2},
					Sku:   "SKU002",
					Name:  "Item 2",
					Price: 50.0,
				}, nil)
			},
			expectedError: service.ErrInvalidAmount,
			checkOrder:    nil,
		},
		{
			name:       "Failure - Database error",
			customerID: "customer123",
			items: []dto.ItemRequest{
				{Sku: "SKU001", Quantity: 2},
			},
			payments: []dto.PaymentRequest{
				{Method: model.COD, Amount: 100.0},
			},
			mockSetup: func(orderRepo *MockOrderRepository, itemRepo *MockItemRepository) {
				// Mock item repository responses
				itemRepo.On("GetBySku", mock.Anything, "SKU001").Return(&model.Item{
					Base:  model.Base{ID: 1},
					Sku:   "SKU001",
					Name:  "Item 1",
					Price: 50.0,
				}, nil)

				// Mock database error on order creation
				orderRepo.On("Create", mock.Anything, mock.AnythingOfType("*model.Order")).Return(errors.New("database error"))
			},
			expectedError: errors.New("failed to create order: database error"),
			checkOrder:    nil,
		},
		{
			name:       "Success - Multiple payment methods",
			customerID: "customer123",
			items: []dto.ItemRequest{
				{Sku: "SKU001", Quantity: 2},
			},
			payments: []dto.PaymentRequest{
				{Method: model.COD, Amount: 50.0},
				{Method: model.VNPAY, Amount: 50.0},
			},
			mockSetup: func(orderRepo *MockOrderRepository, itemRepo *MockItemRepository) {
				// Mock item repository responses
				itemRepo.On("GetBySku", mock.Anything, "SKU001").Return(&model.Item{
					Base:  model.Base{ID: 1},
					Sku:   "SKU001",
					Name:  "Item 1",
					Price: 50.0,
				}, nil)

				// Mock successful order creation
				orderRepo.On("Create", mock.Anything, mock.AnythingOfType("*model.Order")).Run(func(args mock.Arguments) {
					// Set ID when order is created
					order := args.Get(1).(*model.Order)
					order.ID = 1
					order.CreatedAt = time.Now()
					order.UpdatedAt = time.Now()

					// Set order IDs for nested items
					for i := range order.Items {
						order.Items[i].OrderID = order.ID
					}

					// Set order IDs for payments
					for i := range order.Payments {
						order.Payments[i].OrderID = order.ID
					}
				}).Return(nil)
			},
			expectedError: nil,
			checkOrder: func(t *testing.T, order *model.Order) {
				assert.NotNil(t, order)
				assert.Equal(t, int64(1), order.ID)
				assert.Equal(t, "customer123", order.CustomerID)
				assert.Equal(t, 100.0, order.TotalAmount)
				assert.Equal(t, model.OrderPending, order.Status)

				// Check payments
				assert.Len(t, order.Payments, 2)
				assert.Equal(t, model.COD, order.Payments[0].Method)
				assert.Equal(t, 50.0, order.Payments[0].Amount)
				assert.Equal(t, model.VNPAY, order.Payments[1].Method)
				assert.Equal(t, 50.0, order.Payments[1].Amount)
			},
		},
		{
			name:       "Success - Zero items (edge case)",
			customerID: "customer123",
			items:      []dto.ItemRequest{},
			payments:   []dto.PaymentRequest{},
			mockSetup: func(orderRepo *MockOrderRepository, itemRepo *MockItemRepository) {
				// Mock successful order creation with no items
				orderRepo.On("Create", mock.Anything, mock.AnythingOfType("*model.Order")).Run(func(args mock.Arguments) {
					// Set ID when order is created
					order := args.Get(1).(*model.Order)
					order.ID = 1
					order.CreatedAt = time.Now()
					order.UpdatedAt = time.Now()
				}).Return(nil)
			},
			expectedError: nil,
			checkOrder: func(t *testing.T, order *model.Order) {
				assert.NotNil(t, order)
				assert.Equal(t, int64(1), order.ID)
				assert.Equal(t, "customer123", order.CustomerID)
				assert.Equal(t, 0.0, order.TotalAmount)
				assert.Equal(t, model.OrderPending, order.Status)

				// Check items and payments
				assert.Len(t, order.Items, 0)
				assert.Len(t, order.Payments, 0)
			},
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create mock repositories
			mockOrderRepo := new(MockOrderRepository)
			mockItemRepo := new(MockItemRepository)

			// Set up mock expectations
			if tc.mockSetup != nil {
				tc.mockSetup(mockOrderRepo, mockItemRepo)
			}

			// Create the service with mock repositories
			orderService := service.NewOrderService(mockOrderRepo, mockItemRepo)

			// Execute the function being tested
			order, err := orderService.CreateOrder(context.Background(), tc.customerID, tc.items, tc.payments)

			// Check error expectations
			if tc.expectedError != nil {
				assert.Error(t, err)
				if errors.Is(err, tc.expectedError) {
					assert.True(t, errors.Is(err, tc.expectedError))
				} else {
					assert.Contains(t, err.Error(), tc.expectedError.Error())
				}
				assert.Nil(t, order)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, order)

				// Verify order details if a checking function is provided
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
