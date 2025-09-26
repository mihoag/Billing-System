package tests

import (
	"billing-system/billing_service/internal/model"
	"billing-system/billing_service/internal/repository"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestInvoiceRepositoryCreate(t *testing.T) {
	// Test cases for table-driven tests
	testCases := []struct {
		name          string
		invoice       *model.Invoice
		mockSetup     func(mock sqlmock.Sqlmock)
		expectedError error
	}{
		{
			name: "Success - Create invoice with items",
			invoice: &model.Invoice{
				OrderID:     1,
				ShipmentID:  100,
				TotalAmount: 99.99,
				Items: []model.InvoiceItem{
					{
						Quantity: 1,
						ItemID:   1,
					},
				},
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				// Expect transaction begin
				mock.ExpectBegin()

				// Expect invoice creation
				mock.ExpectQuery(`INSERT INTO "invoices"`).
					WithArgs(
						AnyTime(), AnyTime(), nil, // Base fields
						1, 100, 99.99, // Invoice fields (order_id, shipment_id, total_amount)
					).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

				// Expect InvoiceItem creation
				mock.ExpectQuery(`INSERT INTO "invoice_items"`).
					WithArgs(
						AnyTime(), AnyTime(), nil, // Base fields
						1, 1, 1, // InvoiceItem fields (invoice_id, quantity, item_id)
					).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

				// Expect transaction commit
				mock.ExpectCommit()
			},
			expectedError: nil,
		},
		{
			name: "Error - Database error during invoice creation",
			invoice: &model.Invoice{
				OrderID:     2,
				ShipmentID:  200,
				TotalAmount: 199.99,
				Items:       []model.InvoiceItem{},
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				// Expect transaction begin
				mock.ExpectBegin()

				// Expect invoice creation with error
				mock.ExpectQuery(`INSERT INTO "invoices"`).
					WithArgs(
						AnyTime(), AnyTime(), nil, // Base fields
						2, 200, 199.99, // Invoice fields
					).
					WillReturnError(errors.New("database error"))

				// Expect transaction rollback
				mock.ExpectRollback()
			},
			expectedError: errors.New("database error"),
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a new mock database for each test case
			mockDB, err := NewMockDB()
			assert.NoError(t, err)
			defer mockDB.Close()

			// Configure the mock according to the test case
			tc.mockSetup(mockDB.Mock)

			// Create a new invoice repository with the mock database
			invoiceRepo := repository.NewInvoiceRepository(mockDB.DB)

			// Call the method being tested
			err = invoiceRepo.Create(context.Background(), tc.invoice)

			// Check the results
			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}

			// Verify that all expectations were met
			assert.NoError(t, mockDB.ExpectationsWereMet())
		})
	}
}

func TestInvoiceRepositoryGetByOrderID(t *testing.T) {
	// Test cases for table-driven tests
	testCases := []struct {
		name             string
		orderID          int64
		mockSetup        func(mock sqlmock.Sqlmock)
		expectedInvoices []model.Invoice
		expectedError    error
	}{
		{
			name:    "Success - Found invoices for order",
			orderID: 1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				// Invoice rows
				invoiceRows := sqlmock.NewRows(InvoiceColumns()).
					AddRow(1, time.Now(), time.Now(), nil, 1, 100, 99.99).
					AddRow(2, time.Now(), time.Now(), nil, 1, 101, 49.99)

				mock.ExpectQuery(`SELECT (.+) FROM "invoices"`).
					WithArgs(1).
					WillReturnRows(invoiceRows)

				// Items for invoice 1
				itemRows1 := sqlmock.NewRows(InvoiceItemColumns()).
					AddRow(1, time.Now(), time.Now(), nil, 1, 1, 1).
					AddRow(2, time.Now(), time.Now(), nil, 1, 2, 2)

				mock.ExpectQuery(`SELECT (.+) FROM "invoice_items"`).
					WithArgs(1).
					WillReturnRows(itemRows1)

				// Items for invoice 2
				itemRows2 := sqlmock.NewRows(InvoiceItemColumns()).
					AddRow(3, time.Now(), time.Now(), nil, 2, 1, 3)

				mock.ExpectQuery(`SELECT (.+) FROM "invoice_items"`).
					WithArgs(2).
					WillReturnRows(itemRows2)
			},
			expectedInvoices: []model.Invoice{
				{
					Base:        model.Base{ID: 1},
					OrderID:     1,
					ShipmentID:  100,
					TotalAmount: 99.99,
					Items: []model.InvoiceItem{
						{Base: model.Base{ID: 1}, InvoiceID: 1, Quantity: 1, ItemID: 1},
						{Base: model.Base{ID: 2}, InvoiceID: 1, Quantity: 2, ItemID: 2},
					},
				},
				{
					Base:        model.Base{ID: 2},
					OrderID:     1,
					ShipmentID:  101,
					TotalAmount: 49.99,
					Items: []model.InvoiceItem{
						{Base: model.Base{ID: 3}, InvoiceID: 2, Quantity: 1, ItemID: 3},
					},
				},
			},
			expectedError: nil,
		},
		{
			name:    "Success - No invoices found",
			orderID: 2,
			mockSetup: func(mock sqlmock.Sqlmock) {
				// Empty result
				mock.ExpectQuery(`SELECT (.+) FROM "invoices"`).
					WithArgs(2).
					WillReturnRows(sqlmock.NewRows(InvoiceColumns()))
			},
			expectedInvoices: []model.Invoice{},
			expectedError:    nil,
		},
		{
			name:    "Error - Database error",
			orderID: 3,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT (.+) FROM "invoices"`).
					WithArgs(3).
					WillReturnError(errors.New("database error"))
			},
			expectedInvoices: nil,
			expectedError:    errors.New("database error"),
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a new mock database for each test case
			mockDB, err := NewMockDB()
			assert.NoError(t, err)
			defer mockDB.Close()

			// Configure the mock according to the test case
			tc.mockSetup(mockDB.Mock)

			// Create a new invoice repository with the mock database
			invoiceRepo := repository.NewInvoiceRepository(mockDB.DB)

			// Call the method being tested
			invoices, err := invoiceRepo.GetByOrderID(context.Background(), tc.orderID)

			// Check the results
			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedError.Error())
				assert.Nil(t, invoices)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(tc.expectedInvoices), len(invoices))

				// Check each invoice
				for i, expectedInvoice := range tc.expectedInvoices {
					assert.Equal(t, expectedInvoice.ID, invoices[i].ID)
					assert.Equal(t, expectedInvoice.OrderID, invoices[i].OrderID)
					assert.Equal(t, expectedInvoice.ShipmentID, invoices[i].ShipmentID)
					assert.Equal(t, expectedInvoice.TotalAmount, invoices[i].TotalAmount)

					// Check invoice items
					assert.Equal(t, len(expectedInvoice.Items), len(invoices[i].Items))
					for j, expectedItem := range expectedInvoice.Items {
						assert.Equal(t, expectedItem.ID, invoices[i].Items[j].ID)
						assert.Equal(t, expectedItem.InvoiceID, invoices[i].Items[j].InvoiceID)
						assert.Equal(t, expectedItem.Quantity, invoices[i].Items[j].Quantity)
						assert.Equal(t, expectedItem.ItemID, invoices[i].Items[j].ItemID)
					}
				}
			}

			// Verify that all expectations were met
			assert.NoError(t, mockDB.ExpectationsWereMet())
		})
	}
}
