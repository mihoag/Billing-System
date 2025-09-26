package handler

import (
	"billing-system/shipment_service/internal/service"
	"billing-system/shipment_service/pkg/utils"
	pb "billing-system/shipment_service/proto"
	"context"
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
	// Convert proto ShipmentItemRequests to DTO ShipmentItemRequests
	items := utils.ConvertProtoItemsToDTO(req.Items)

	// Call the service layer to create the shipment
	shipment, err := h.shipmentService.CreateShipment(ctx, req.OrderId, items)
	if err != nil {
		return &pb.CreateShipmentResponse{
			Code:    0, // Error code
			Message: err.Error(),
		}, nil
	}

	// Convert the domain shipment to proto shipment data
	shipmentData := utils.ConvertShipmentToProtoData(shipment)

	return &pb.CreateShipmentResponse{
		Code:    1, // Success code
		Message: "Create shipment successfully",
		Data:    shipmentData,
	}, nil
}
