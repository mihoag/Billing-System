package repository

import (
	"billing-system/billing_service/internal/model"
	"context"
	"errors"

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

// GetByID retrieves an item by its ID
func (r *ItemRepositoryImpl) GetByID(ctx context.Context, id int64) (*model.Item, error) {
	if id <= 0 {
		return nil, errors.New("invalid item ID")
	}

	var item model.Item
	result := r.db.WithContext(ctx).First(&item, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("item not found")
		}
		return nil, result.Error
	}

	return &item, nil
}
