package billing

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"billing-system/bff/config"
	billingPb "billing-system/billing_service/proto"
)

type Handler struct {
	BillingConnection *BillingConnectionAdapter
	logger            *zap.Logger
}

func NewHandler() *Handler {
	return &Handler{
		BillingConnection: &BillingConnectionAdapter{},
		logger:            config.Service.Logger,
	}
}

// CreateOrder godoc
// @Summary Create a new order
// @Description Create a new order with items and payment details
// @Tags orders
// @Accept json
// @Produce json
// @Param order body CreateOrderRequest true "Order details"
// @Success 200 {object} CreateOrderResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /orders [post]
func (h *Handler) CreateOrder(ctx *gin.Context) {
	var request CreateOrderRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get billing service client
	client, _, err := h.BillingConnection.NewClient()
	if err != nil {
		h.logger.Error("Failed to create billing client", zap.Error(err))
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
			ItemId:   item.ItemID,
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
		h.logger.Error("Failed to create order", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert response
	response := convertPbOrderToResponse(pbResponse.Order)
	ctx.JSON(http.StatusOK, response)
}

// GetOrder godoc
// @Summary Get an order by ID
// @Description Get detailed information about an order by its ID
// @Tags orders
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Success 200 {object} OrderResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /orders/{id} [get]
func (h *Handler) GetOrder(ctx *gin.Context) {
	// Parse order ID from path
	orderIDStr := ctx.Param("id")
	orderID, err := strconv.ParseInt(orderIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	// Get billing service client
	client, _, err := h.BillingConnection.NewClient()
	if err != nil {
		h.logger.Error("Failed to create billing client", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to billing service"})
		return
	}
	billingClient := client.(billingPb.BillingServiceClient)

	// Call billing service
	pbResponse, err := billingClient.GetOrder(context.Background(), &billingPb.GetOrderRequest{
		OrderId: orderID,
	})
	if err != nil {
		h.logger.Error("Failed to get order", zap.Error(err))
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
