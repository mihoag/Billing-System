package mocks

import (
	"billing-system/billing_service/internal/model"
	"context"

	"github.com/stretchr/testify/mock"
)

// MockInvoiceRepository is a mock implementation of repository.InvoiceRepository
type MockInvoiceRepository struct {
	mock.Mock
}

func (m *MockInvoiceRepository) Create(ctx context.Context, invoice *model.Invoice) error {
	args := m.Called(ctx, invoice)
	return args.Error(0)
}

func (m *MockInvoiceRepository) GetByOrderID(ctx context.Context, orderID int64) ([]model.Invoice, error) {
	args := m.Called(ctx, orderID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Invoice), args.Error(1)
}