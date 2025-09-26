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
	"gorm.io/gorm"
)

func TestItemRepositoryGetBySku(t *testing.T) {
	// Test cases for table-driven tests
	testCases := []struct {
		name          string
		sku           string
		mockSetup     func(mock sqlmock.Sqlmock)
		expectedItem  *model.Item
		expectedError error
	}{
		{
			name: "Success - Item found",
			sku:  "SKU123",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows(ItemColumns()).
					AddRow(1, time.Now(), time.Now(), nil, "Test Item", "SKU123", 99.99)
				mock.ExpectQuery(`SELECT (.+) FROM "items"`).
					WithArgs("SKU123", 1). // GORM adds LIMIT 1 for First()
					WillReturnRows(rows)
			},
			expectedItem: &model.Item{
				Base: model.Base{
					ID: 1,
				},
				Name:  "Test Item",
				Sku:   "SKU123",
				Price: 99.99,
			},
			expectedError: nil,
		},
		{
			name: "Error - Item not found",
			sku:  "NONEXISTENT",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT (.+) FROM "items"`).
					WithArgs("NONEXISTENT", 1). // GORM adds LIMIT 1 for First()
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectedItem:  nil,
			expectedError: gorm.ErrRecordNotFound,
		},
		{
			name: "Error - Database error",
			sku:  "SKU456",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT (.+) FROM "items"`).
					WithArgs("SKU456", 1). // GORM adds LIMIT 1 for First()
					WillReturnError(errors.New("database connection error"))
			},
			expectedItem:  nil,
			expectedError: errors.New("database connection error"),
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

			// Create a new item repository with the mock database
			itemRepo := repository.NewItemRepository(mockDB.DB)

			// Call the method being tested
			item, err := itemRepo.GetBySku(context.Background(), tc.sku)

			// Check the results
			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError.Error(), err.Error())
				assert.Nil(t, item)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, item)
				assert.Equal(t, tc.expectedItem.ID, item.ID)
				assert.Equal(t, tc.expectedItem.Name, item.Name)
				assert.Equal(t, tc.expectedItem.Sku, item.Sku)
				assert.Equal(t, tc.expectedItem.Price, item.Price)
			}

			// Verify that all expectations were met
			assert.NoError(t, mockDB.ExpectationsWereMet())
		})
	}
}
