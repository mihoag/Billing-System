package shipment

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	shipmentPb "billing-system/shipment_service/proto"
)

// Handler struct
type Handler struct {
	ShipmentConnection *ShipmentConnectionAdapter
}

// NewHandler creates a new shipment handler
func NewHandler() *Handler {
	return &Handler{
		ShipmentConnection: &ShipmentConnectionAdapter{},
	}
}

// CreateShipmentHandler handles HTTP request to create a shipment
func (h *Handler) CreateShipment(ctx *gin.Context) {
	var request CreateShipmentRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get shipment service client
	client, _, err := h.ShipmentConnection.NewClient()
	if err != nil {
		log.Println("Error connecting to shipment service:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to shipment service"})
		return
	}

	shipmentClient := client.(shipmentPb.ShipmentServiceClient)

	// Convert request to protobuf
	protoItems := make([]*shipmentPb.ShipmentItemRequest, len(request.Items))
	for i, item := range request.Items {
		protoItems[i] = &shipmentPb.ShipmentItemRequest{
			Sku:      item.Sku,
			Quantity: item.Quantity,
		}
	}

	protoReq := &shipmentPb.CreateShipmentRequest{
		OrderId: request.OrderID,
		Items:   protoItems,
	}

	// Call shipment service
	protoResp, err := shipmentClient.CreateShipment(ctx, protoReq)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert response
	response := &ShipmentResponse{
		Code:    int(protoResp.Code),
		Message: protoResp.Message,
		Data:    protoResp.Data,
	}

	ctx.JSON(http.StatusOK, response)
}
