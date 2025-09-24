package billing

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	billingPb "billing-system/billing_service/proto"
)

type Handler struct {
	BillingConnection *BillingConnectionAdapter
}

func NewHandler() *Handler {
	return &Handler{
		BillingConnection: &BillingConnectionAdapter{},
	}
}

func (h *Handler) CreateOrder(ctx *gin.Context) {
	var request CreateOrderRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get billing service client
	client, _, err := h.BillingConnection.NewClient()
	if err != nil {
		log.Println("Error connecting to billing service:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to billing service"})
		return
	}
	billingClient := client.(billingPb.BillingServiceClient)

	// Convert request to protobuf
	pbRequest := &billingPb.CreateOrderRequest{
		CustomerId: request.CustomerID,
		Items:      make([]*billingPb.ItemRequest, len(request.Items)),
		Payments:   make([]*billingPb.PaymentRequest, len(request.Payments)),
	}

	// Convert items
	for i, item := range request.Items {
		pbRequest.Items[i] = &billingPb.ItemRequest{
			Sku:      item.Sku,
			Quantity: int32(item.Quantity),
			Price:    item.Price,
		}
	}

	// Convert payments
	for i, payment := range request.Payments {
		pbRequest.Payments[i] = &billingPb.PaymentRequest{
			Method: payment.Method,
			Amount: payment.Amount,
		}
	}

	// Call billing service
	pbResponse, err := billingClient.CreateOrder(ctx, pbRequest)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert response
	response := convertPbOrderToResponse(pbResponse.Order)
	ctx.JSON(http.StatusOK, response)
}

// Helper functions for request/response conversion
func convertPbOrderToResponse(pbOrder *billingPb.Order) OrderResponse {
	response := OrderResponse{
		ID:          pbOrder.Id,
		CustomerID:  pbOrder.CustomerId,
		TotalAmount: pbOrder.TotalAmount,
		Status:      pbOrder.Status.String(),
		CreatedAt:   pbOrder.CreatedAt,
		UpdatedAt:   pbOrder.UpdatedAt,
		Items:       make([]OrderItemResponse, len(pbOrder.Items)),
		Payments:    make([]PaymentResponse, len(pbOrder.Payments)),
	}

	// Convert items
	for i, item := range pbOrder.Items {
		response.Items[i] = OrderItemResponse{
			ID:       item.Id,
			OrderID:  item.OrderId,
			ItemID:   item.ItemId,
			Quantity: int(item.Quantity),
		}
	}

	// Convert payments
	for i, payment := range pbOrder.Payments {
		response.Payments[i] = PaymentResponse{
			ID:      payment.Id,
			OrderID: payment.OrderId,
			Method:  payment.Method,
			Amount:  payment.Amount,
		}
	}

	return response
}
