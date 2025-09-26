package utils

import (
	"billing-system/shipment_service/internal/dto"
	"billing-system/shipment_service/internal/model"
	pb "billing-system/shipment_service/proto"
	"time"
)

// ConvertProtoItemsToDTO converts proto ShipmentItemRequests to DTO ShipmentItemRequests
func ConvertProtoItemsToDTO(protoItems []*pb.ShipmentItemRequest) []dto.ShipmentItemRequest {
	items := make([]dto.ShipmentItemRequest, len(protoItems))
	for i, item := range protoItems {
		items[i] = dto.ShipmentItemRequest{
			Sku:      item.Sku,
			Quantity: int(item.Quantity),
		}
	}
	return items
}

// ConvertShipmentToProtoData converts a domain Shipment to proto ShipmentData
func ConvertShipmentToProtoData(shipment *model.Shipment) *pb.ShipmentData {
	shipmentData := &pb.ShipmentData{
		ShipmentId: shipment.ID,
		OrderId:    shipment.OrderID,
		Status:     string(shipment.Status),
		CreatedAt:  shipment.CreatedAt.Format(time.RFC3339),
	}

	// Convert shipment items
	protoItems := make([]*pb.ShipmentItem, len(shipment.Items))
	for i, item := range shipment.Items {
		protoItems[i] = &pb.ShipmentItem{
			Sku:      item.Sku,
			Quantity: int32(item.Quantity),
		}
	}
	shipmentData.Items = protoItems

	return shipmentData
}
