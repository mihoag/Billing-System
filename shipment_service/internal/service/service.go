package service

import (
	"billing-system/shipment_service/internal/dto"
	"billing-system/shipment_service/internal/model"
	"context"
)

type ShipmentService interface {
	CreateShipment(ctx context.Context, orderID int64, items []dto.ShipmentItemRequest) (*model.Shipment, error)
}
