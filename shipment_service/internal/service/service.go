package service

import (
	"billing-system/shipment_service/internal/model"
	"context"
)

type ShipmentService interface {
	CreateShipment(ctx context.Context, customerID string, items []ShipmentItemRequest) (*model.Shipment, error)
}
