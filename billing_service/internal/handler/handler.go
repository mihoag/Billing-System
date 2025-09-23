package billing_handler

import (
	"billing-system/billing_service/internal/model"
	"billing-system/billing_service/internal/service"
	pb "billing-system/billing_service/proto"
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// OrderHandler handles gRPC requestsW related to orders
type OrderHandler struct {
	pb.UnimplementedBillingServiceServer
	orderService service.OrderService
}

// NewOrderHandler creates a new OrderHandler
func NewOrderHandler(orderService service.OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}

// CreateOrder handles the gRPC request to create a new order
func (h *OrderHandler) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	// Convert proto ItemRequests to service ItemRequests
	items := make([]service.ItemRequest, len(req.Items))
	for i, item := range req.Items {
		items[i] = service.ItemRequest{
			ItemID:   item.ItemId,
			Quantity: int(item.Quantity),
			Price:    item.Price,
		}
	}

	// Convert proto PaymentRequests to service PaymentRequests
	payments := make([]service.PaymentRequest, len(req.Payments))
	for i, payment := range req.Payments {
		payments[i] = service.PaymentRequest{
			Method: model.PaymentMethod(payment.Method),
			Amount: payment.Amount,
		}
	}

	// Call the service layer to create the order
	order, err := h.orderService.CreateOrder(ctx, req.CustomerId, items, payments)
	if err != nil {
		return nil, mapErrorToGRPCStatus(err).Err()
	}

	// Convert the domain order to proto order
	protoOrder := convertOrderToProto(order)

	return &pb.CreateOrderResponse{
		Order: protoOrder,
	}, nil
}

// convertOrderToProto converts a domain Order to a proto Order
func convertOrderToProto(order *model.Order) *pb.Order {
	protoOrder := &pb.Order{
		Id:          order.ID,
		CustomerId:  order.CustomerID,
		TotalAmount: order.TotalAmount,
		Status:      mapOrderStatus(order.Status),
		CreatedAt:   order.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   order.UpdatedAt.Format(time.RFC3339),
	}

	// Convert order items
	protoItems := make([]*pb.OrderItem, len(order.Items))
	for i, item := range order.Items {
		protoItems[i] = &pb.OrderItem{
			Id:       item.ID,
			OrderId:  item.OrderID,
			ItemId:   item.ItemId,
			Quantity: int32(item.Quantity),
		}
	}
	protoOrder.Items = protoItems

	// Convert payments
	protoPayments := make([]*pb.Payment, len(order.Payments))
	for i, payment := range order.Payments {
		protoPayments[i] = &pb.Payment{
			Id:      payment.ID,
			OrderId: payment.OrderID,
			Method:  string(payment.Method),
			Amount:  payment.Amount,
		}
	}
	protoOrder.Payments = protoPayments

	return protoOrder
}

// mapOrderStatus maps a domain OrderStatus to a proto OrderStatus
func mapOrderStatus(status model.OrderStatus) pb.OrderStatus {
	switch status {
	case model.Pending:
		return pb.OrderStatus_PENDING
	case model.Success:
		return pb.OrderStatus_SUCCESS
	case model.Failed:
		return pb.OrderStatus_FAILED
	default:
		return pb.OrderStatus_PENDING
	}
}

// mapErrorToGRPCStatus maps service errors to gRPC status errors
func mapErrorToGRPCStatus(err error) *status.Status {
	switch err {
	case service.ErrItemNotFound, service.ErrOrderNotFound:
		return status.New(codes.NotFound, err.Error())
	case service.ErrInvalidQuantity, service.ErrInvalidAmount, service.ErrInsufficientPayment:
		return status.New(codes.InvalidArgument, err.Error())
	default:
		return status.New(codes.Internal, "internal server error")
	}
}
