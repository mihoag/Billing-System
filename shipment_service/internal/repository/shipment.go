package repository

import (
	"billing-system/shipment_service/internal/model"
	"context"

	"gorm.io/gorm"
)

type shipmentRepository struct {
	db *gorm.DB
}

// NewShipmentRepository creates a new shipment repository
func NewShipmentRepository(db *gorm.DB) ShipmentRepository {
	return &shipmentRepository{
		db: db,
	}
}

// Create creates a new shipment with its items
func (r *shipmentRepository) Create(ctx context.Context, shipment *model.Shipment) error {
	return r.db.WithContext(ctx).Create(shipment).Error
}

// Update updates an existing shipment
func (r *shipmentRepository) Update(ctx context.Context, shipment *model.Shipment) error {
	return r.db.WithContext(ctx).Save(shipment).Error
}
