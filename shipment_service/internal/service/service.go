package service

import (
	"billing-system/shipment_service/internal/model"
	"context"
)

type ShipmentService interface {
	CreateShipment(ctx context.Context, orderID int64, items []ShipmentItemRequest) (*model.Shipment, error)
}

type ShipmentItemRequest struct {
	Sku      string
	Quantity int
}
