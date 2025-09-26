package tests

import (
	"billing-system/billing_service/internal/model"
	"billing-system/billing_service/internal/repository"
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestOrderRepositoryCreate(t *testing.T) {
	// Test cases for table-driven tests
	testCases := []struct {
		name          string
		order         *model.Order
		mockSetup     func(mock sqlmock.Sqlmock)
		expectedError error
	}{
		{
			name: "Success - Create order with items and payments",
			order: &model.Order{
				CustomerID:  "CUST123",
				TotalAmount: 199.98,
				Status:      model.OrderPending,
				Items: []model.OrderItem{
					{
						Quantity: 2,
						ItemID:   1,
					},
				},
				Payments: []model.Payment{
					{
						Method: model.COD,
						Amount: 199.98,
					},
				},
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				// Expect transaction begin
				mock.ExpectBegin()

				// Expect order creation
				mock.ExpectQuery(`INSERT INTO "orders"`).
					WithArgs(
						AnyTime(), AnyTime(), nil, // Base fields
						"CUST123", 199.98, model.OrderPending, // Order fields
					).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

				// Expect OrderItem creation
				mock.ExpectQuery(`INSERT INTO "order_items"`).
					WithArgs(
						AnyTime(), AnyTime(), nil, // Base fields
						1, 2, 1, // OrderItem fields (order_id, quantity, item_id)
					).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

				// Expect Payment creation
				mock.ExpectQuery(`INSERT INTO "payments"`).
					WithArgs(
						AnyTime(), AnyTime(), nil, // Base fields
						1, model.COD, 199.98, // Payment fields (order_id, method, amount)
					).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

				// Expect transaction commit
				mock.ExpectCommit()
			},
			expectedError: nil,
		},
		{
			name: "Error - Database error during order creation",
			order: &model.Order{
				CustomerID:  "CUST456",
				TotalAmount: 99.99,
				Status:      model.OrderPending,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				// Expect transaction begin
				mock.ExpectBegin()

				// Expect order creation with error
				mock.ExpectQuery(`INSERT INTO "orders"`).
					WithArgs(
						AnyTime(), AnyTime(), nil, // Base fields
						"CUST456", 99.99, model.OrderPending, // Order fields
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

			// Create a new order repository with the mock database
			orderRepo := repository.NewOrderRepository(mockDB.DB)

			// Call the method being tested
			err = orderRepo.Create(context.Background(), tc.order)

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
