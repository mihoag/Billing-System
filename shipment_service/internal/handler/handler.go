package handler

import (
	"billing-system/shipment_service/internal/dto"
	"billing-system/shipment_service/internal/model"
	"billing-system/shipment_service/internal/service"
	pb "billing-system/shipment_service/proto"
	"context"
	"time"
)

// ShipmentHandler handles gRPC requests related to shipments
type ShipmentHandler struct {
	pb.UnimplementedShipmentServiceServer
	shipmentService service.ShipmentService
}

// NewShipmentHandler creates a new ShipmentHandler
func NewShipmentHandler(shipmentService service.ShipmentService) *ShipmentHandler {
	return &ShipmentHandler{
		shipmentService: shipmentService,
	}
}

// CreateShipment handles the gRPC request to create a new shipment
func (h *ShipmentHandler) CreateShipment(ctx context.Context, req *pb.CreateShipmentRequest) (*pb.CreateShipmentResponse, error) {
	// Convert proto ShipmentItemRequests to service ShipmentItemRequests
	items := make([]dto.ShipmentItemRequest, len(req.Items))
	for i, item := range req.Items {
		items[i] = dto.ShipmentItemRequest{
			Sku:      item.Sku,
			Quantity: int(item.Quantity),
		}
	}

	// Call the service layer to create the shipment
	shipment, err := h.shipmentService.CreateShipment(ctx, req.OrderId, items)
	if err != nil {
		return &pb.CreateShipmentResponse{
			Code:    0, // Error code
			Message: err.Error(),
		}, nil
	}

	// Convert the domain shipment to proto shipment data
	shipmentData := convertShipmentToProtoData(shipment)

	return &pb.CreateShipmentResponse{
		Code:    1, // Success code
		Message: "Create shipment successfully",
		Data:    shipmentData,
	}, nil
}

// convertShipmentToProtoData converts a domain Shipment to proto ShipmentData
func convertShipmentToProtoData(shipment *model.Shipment) *pb.ShipmentData {
	shipmentData := &pb.ShipmentData{
		ShipmentId: shipment.ID,
		OrderId:    shipment.ID,
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
