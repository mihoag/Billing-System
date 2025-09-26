package mocks

import (
	"billing-system/billing_service/internal/model"
	"context"

	"github.com/stretchr/testify/mock"
)

// MockOrderRepository is a mock implementation of repository.OrderRepository
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

// MockItemRepository is a mock implementation of repository.ItemRepository
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