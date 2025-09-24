package service

import (
	"billing-system/shipment_service/internal/model"
	"billing-system/shipment_service/internal/repository"
	"context"
	"fmt"
)

type ShipmentServiceImpl struct {
	shipmentRepo repository.ShipmentRepository
}

func NewShipmentService(shipmentRepo repository.ShipmentRepository) ShipmentService {
	return &ShipmentServiceImpl{
		shipmentRepo: shipmentRepo,
	}
}

func (s *ShipmentServiceImpl) CreateShipment(ctx context.Context, orderID int64, items []ShipmentItemRequest) (*model.Shipment, error) {

	if len(items) == 0 {
		return nil, fmt.Errorf("at least one item is required")
	}

	// Validate items
	var shipmentItems []model.ShipmentItem
	for _, itemReq := range items {
		if itemReq.Sku == "" {
			return nil, fmt.Errorf("SKU is required for all items")
		}
		if itemReq.Quantity <= 0 {
			return nil, fmt.Errorf("quantity must be greater than 0 for SKU %s", itemReq.Sku)
		}

		shipmentItems = append(shipmentItems, model.ShipmentItem{
			Sku:      itemReq.Sku,
			Quantity: itemReq.Quantity,
		})
	}

	// Create shipment
	shipment := &model.Shipment{
		OrderID: orderID,
		Status:  model.Confirmed,
		Items:   shipmentItems,
	}

	if err := s.shipmentRepo.Create(ctx, shipment); err != nil {
		return nil, fmt.Errorf("failed to create shipment: %w", err)
	}

	if err := s.shipmentRepo.Update(ctx, shipment); err != nil {
		return nil, fmt.Errorf("failed to update shipment status: %w", err)
	}

	return shipment, nil
}
