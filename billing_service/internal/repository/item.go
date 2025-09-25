package repository

import (
	"billing-system/billing_service/internal/model"
	"context"

	"gorm.io/gorm"
)

// ItemRepositoryImpl implements the ItemRepository interface
type ItemRepositoryImpl struct {
	db *gorm.DB
}

// NewItemRepository creates a new instance of ItemRepositoryImpl
func NewItemRepository(db *gorm.DB) ItemRepository {
	return &ItemRepositoryImpl{
		db: db,
	}
}

// GetByID retrieves an item by its sku
func (r *ItemRepositoryImpl) GetBySku(ctx context.Context, sku string) (*model.Item, error) {
	var item model.Item
	err := r.db.WithContext(ctx).Where("sku = ?", sku).First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}
