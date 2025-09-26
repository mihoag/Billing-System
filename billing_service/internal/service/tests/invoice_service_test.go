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

func TestInvoiceService_CreateInvoice(t *testing.T) {
	// Define test time for consistent timestamps
	testTime := time.Now()

	// Define test cases
	testCases := []struct {
		name          string
		shipmentID    int64
		orderID       int64
		itemRequests  []dto.InvoiceItemRequest
		mockSetup     func(*mocks.MockInvoiceRepository, *mocks.MockOrderRepository, *mocks.MockItemRepository)
		expectedError string
		checkInvoice  func(*testing.T, *model.Invoice)
	}{
		{
			name:       "Success - Create invoice for all items in order",
			shipmentID: 101,
			orderID:    1,
			itemRequests: []dto.InvoiceItemRequest{
				{Sku: "SKU001", Quantity: 2},
				{Sku: "SKU002", Quantity: 1},
			},
			mockSetup: func(invoiceRepo *mocks.MockInvoiceRepository, orderRepo *mocks.MockOrderRepository, itemRepo *mocks.MockItemRepository) {
				// Mock order repository to return an order with items
				orderRepo.On("GetByID", mock.Anything, int64(1)).Return(&model.Order{
					Base: model.Base{ID: 1},
					Items: []model.OrderItem{
						{
							ItemID:   1,
							Quantity: 2,
							Item: model.Item{
								Base:  model.Base{ID: 1},
								Name:  "Item 1",
								Sku:   "SKU001",
								Price: 100,
							},
						},
						{
							ItemID:   2,
							Quantity: 1,
							Item: model.Item{
								Base:  model.Base{ID: 2},
								Name:  "Item 2",
								Sku:   "SKU002",
								Price: 100,
							},
						},
					},
				}, nil)

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

				// Mock existing invoices for this order (none in this case)
				invoiceRepo.On("GetByOrderID", mock.Anything, int64(1)).Return([]model.Invoice{}, nil)

				// Mock invoice creation
				invoiceRepo.On("Create", mock.Anything, mock.MatchedBy(func(invoice *model.Invoice) bool {
					return invoice.OrderID == 1 &&
						invoice.ShipmentID == 101 &&
						invoice.TotalAmount == 300.0 &&
						len(invoice.Items) == 2
				})).Return(nil).Run(func(args mock.Arguments) {
					// Set ID and timestamps when the invoice is created
					invoice := args.Get(1).(*model.Invoice)
					invoice.ID = 1001
					invoice.CreatedAt = testTime
					invoice.UpdatedAt = testTime

					// Set IDs for invoice items
					for i := range invoice.Items {
						invoice.Items[i].ID = int64(i + 1)
						invoice.Items[i].InvoiceID = 1001
					}
				})
			},
			expectedError: "",
			checkInvoice: func(t *testing.T, invoice *model.Invoice) {
				assert.Equal(t, int64(1001), invoice.ID)
				assert.Equal(t, int64(1), invoice.OrderID)
				assert.Equal(t, int64(101), invoice.ShipmentID)
				assert.Equal(t, 300.0, invoice.TotalAmount)
				assert.Len(t, invoice.Items, 2)

				// Check invoice items
				assert.Equal(t, 2, invoice.Items[0].Quantity)
				assert.Equal(t, int64(1), invoice.Items[0].ItemID)
				assert.Equal(t, 1, invoice.Items[1].Quantity)
				assert.Equal(t, int64(2), invoice.Items[1].ItemID)
			},
		},
		{
			name:       "Success - Create partial invoice (partial fulfillment)",
			shipmentID: 102,
			orderID:    1,
			itemRequests: []dto.InvoiceItemRequest{
				{Sku: "SKU001", Quantity: 1}, // Only partial quantity
			},
			mockSetup: func(invoiceRepo *mocks.MockInvoiceRepository, orderRepo *mocks.MockOrderRepository, itemRepo *mocks.MockItemRepository) {
				// Mock order repository to return an order with items
				orderRepo.On("GetByID", mock.Anything, int64(1)).Return(&model.Order{
					Base: model.Base{ID: 1},
					Items: []model.OrderItem{
						{
							ItemID:   1,
							Quantity: 2,
							Item: model.Item{
								Base:  model.Base{ID: 1},
								Name:  "Item 1",
								Sku:   "SKU001",
								Price: 100,
							},
						},
						{
							ItemID:   2,
							Quantity: 1,
							Item: model.Item{
								Base:  model.Base{ID: 2},
								Name:  "Item 2",
								Sku:   "SKU002",
								Price: 100,
							},
						},
					},
				}, nil)

				// Mock item repository response
				itemRepo.On("GetBySku", mock.Anything, "SKU001").Return(&model.Item{
					Base:  model.Base{ID: 1},
					Name:  "Item 1",
					Sku:   "SKU001",
					Price: 100,
				}, nil)

				// Mock existing invoices for this order (none)
				invoiceRepo.On("GetByOrderID", mock.Anything, int64(1)).Return([]model.Invoice{}, nil)

				// Mock invoice creation
				invoiceRepo.On("Create", mock.Anything, mock.MatchedBy(func(invoice *model.Invoice) bool {
					return invoice.OrderID == 1 &&
						invoice.ShipmentID == 102 &&
						invoice.TotalAmount == 100.0 &&
						len(invoice.Items) == 1
				})).Return(nil).Run(func(args mock.Arguments) {
					invoice := args.Get(1).(*model.Invoice)
					invoice.ID = 1002
					invoice.CreatedAt = testTime
					invoice.UpdatedAt = testTime

					for i := range invoice.Items {
						invoice.Items[i].ID = int64(i + 1)
						invoice.Items[i].InvoiceID = 1002
					}
				})
			},
			expectedError: "",
			checkInvoice: func(t *testing.T, invoice *model.Invoice) {
				assert.Equal(t, int64(1002), invoice.ID)
				assert.Equal(t, int64(1), invoice.OrderID)
				assert.Equal(t, int64(102), invoice.ShipmentID)
				assert.Equal(t, 100.0, invoice.TotalAmount)
				assert.Len(t, invoice.Items, 1)
				assert.Equal(t, 1, invoice.Items[0].Quantity)
				assert.Equal(t, int64(1), invoice.Items[0].ItemID)
			},
		},
		{
			name:       "Error - Order not found",
			shipmentID: 103,
			orderID:    999, // Non-existent order
			itemRequests: []dto.InvoiceItemRequest{
				{Sku: "SKU001", Quantity: 1},
			},
			mockSetup: func(invoiceRepo *mocks.MockInvoiceRepository, orderRepo *mocks.MockOrderRepository, itemRepo *mocks.MockItemRepository) {
				// Mock order repository to return error for non-existent order
				orderRepo.On("GetByID", mock.Anything, int64(999)).Return(nil, errors.New("order not found"))
			},
			expectedError: "order not found",
			checkInvoice:  nil,
		},
		{
			name:       "Error - Quantity exceeds available (after previous invoices)",
			shipmentID: 107,
			orderID:    1,
			itemRequests: []dto.InvoiceItemRequest{
				{Sku: "SKU001", Quantity: 2}, // Try to invoice 2 more when only 1 remaining
			},
			mockSetup: func(invoiceRepo *mocks.MockInvoiceRepository, orderRepo *mocks.MockOrderRepository, itemRepo *mocks.MockItemRepository) {
				// Mock order repository to return an order with items
				orderRepo.On("GetByID", mock.Anything, int64(1)).Return(&model.Order{
					Base: model.Base{ID: 1},
					Items: []model.OrderItem{
						{
							ItemID:   1,
							Quantity: 2,
							Item: model.Item{
								Base:  model.Base{ID: 1},
								Name:  "Item 1",
								Sku:   "SKU001",
								Price: 100,
							},
						},
					},
				}, nil)

				// Mock item repository
				itemRepo.On("GetBySku", mock.Anything, "SKU001").Return(&model.Item{
					Base:  model.Base{ID: 1},
					Name:  "Item 1",
					Sku:   "SKU001",
					Price: 100,
				}, nil)

				// Mock existing invoices with 1 item already consumed
				invoiceRepo.On("GetByOrderID", mock.Anything, int64(1)).Return([]model.Invoice{
					{
						Base:    model.Base{ID: 1000},
						OrderID: 1,
						Items: []model.InvoiceItem{
							{
								ItemID:   1,
								Quantity: 1, // Already invoiced 1 out of 2
							},
						},
					},
				}, nil)
			},
			expectedError: "exceeds available quantity",
			checkInvoice:  nil,
		},
		{
			name:       "Error - Database error when creating invoice",
			shipmentID: 108,
			orderID:    1,
			itemRequests: []dto.InvoiceItemRequest{
				{Sku: "SKU001", Quantity: 1},
			},
			mockSetup: func(invoiceRepo *mocks.MockInvoiceRepository, orderRepo *mocks.MockOrderRepository, itemRepo *mocks.MockItemRepository) {
				// Mock order repository to return an order with items
				orderRepo.On("GetByID", mock.Anything, int64(1)).Return(&model.Order{
					Base: model.Base{ID: 1},
					Items: []model.OrderItem{
						{
							ItemID:   1,
							Quantity: 2,
							Item: model.Item{
								Base:  model.Base{ID: 1},
								Name:  "Item 1",
								Sku:   "SKU001",
								Price: 100,
							},
						},
					},
				}, nil)

				// Mock item repository
				itemRepo.On("GetBySku", mock.Anything, "SKU001").Return(&model.Item{
					Base:  model.Base{ID: 1},
					Name:  "Item 1",
					Sku:   "SKU001",
					Price: 100,
				}, nil)

				// Mock existing invoices (none)
				invoiceRepo.On("GetByOrderID", mock.Anything, int64(1)).Return([]model.Invoice{}, nil)

				// Mock database error during invoice creation
				invoiceRepo.On("Create", mock.Anything, mock.AnythingOfType("*model.Invoice")).Return(errors.New("database connection error"))
			},
			expectedError: "database connection error",
			checkInvoice:  nil,
		},
		{
			name:       "Error - Failed to retrieve existing invoices",
			shipmentID: 109,
			orderID:    1,
			itemRequests: []dto.InvoiceItemRequest{
				{Sku: "SKU001", Quantity: 1},
			},
			mockSetup: func(invoiceRepo *mocks.MockInvoiceRepository, orderRepo *mocks.MockOrderRepository, itemRepo *mocks.MockItemRepository) {
				// Mock order repository to return an order with items
				orderRepo.On("GetByID", mock.Anything, int64(1)).Return(&model.Order{
					Base: model.Base{ID: 1},
					Items: []model.OrderItem{
						{
							ItemID:   1,
							Quantity: 2,
							Item: model.Item{
								Base:  model.Base{ID: 1},
								Name:  "Item 1",
								Sku:   "SKU001",
								Price: 100,
							},
						},
					},
				}, nil)

				// Mock error when retrieving existing invoices
				invoiceRepo.On("GetByOrderID", mock.Anything, int64(1)).Return(nil, errors.New("failed to retrieve existing invoices"))
			},
			expectedError: "failed to retrieve existing invoices",
			checkInvoice:  nil,
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create mocks
			mockInvoiceRepo := new(mocks.MockInvoiceRepository)
			mockOrderRepo := new(mocks.MockOrderRepository)
			mockItemRepo := new(mocks.MockItemRepository)

			// Set up mocks
			tc.mockSetup(mockInvoiceRepo, mockOrderRepo, mockItemRepo)

			// Create service with mocks
			invoiceService := service.NewInvoiceService(mockInvoiceRepo, mockOrderRepo, mockItemRepo)

			// Call the method being tested
			invoice, err := invoiceService.CreateInvoice(context.Background(), tc.shipmentID, tc.orderID, tc.itemRequests)

			// Check errors
			if tc.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedError)
				assert.Nil(t, invoice)
			} else {
				assert.NoError(t, err)
				require.NotNil(t, invoice)
				if tc.checkInvoice != nil {
					tc.checkInvoice(t, invoice)
				}
			}

			// Verify that all expected mock calls were made
			mockInvoiceRepo.AssertExpectations(t)
			mockOrderRepo.AssertExpectations(t)
			mockItemRepo.AssertExpectations(t)
		})
	}
}
