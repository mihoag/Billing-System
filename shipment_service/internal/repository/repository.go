package repository

import (
	"billing-system/shipment_service/internal/model"
	"context"
)

type ShipmentRepository interface {
	Create(ctx context.Context, shipment *model.Shipment) error
	Update(ctx context.Context, shipment *model.Shipment) error
}
