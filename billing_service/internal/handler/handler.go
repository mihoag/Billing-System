package billing_handler

import (
	"billing-system/billing_service/internal/service"
	"billing-system/billing_service/pkg/utils"
	pb "billing-system/billing_service/proto"
	"context"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// OrderHandler handles gRPC requestsW related to orders
type OrderHandler struct {
	pb.UnimplementedBillingServiceServer
	orderService   service.OrderService
	invoiceService service.InvoiceService
}

// NewOrderHandler creates a new OrderHandler
func NewOrderHandler(orderService service.OrderService, invoiceService service.InvoiceService) *OrderHandler {
	return &OrderHandler{
		orderService:   orderService,
		invoiceService: invoiceService,
	}
}

// CreateOrder handles the gRPC request to create a new order
func (h *OrderHandler) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	// Convert proto requests to DTOs
	items := utils.ProtoItemRequestsToDTO(req.Items)
	payments := utils.ProtoPaymentRequestsToDTO(req.Payments)

	// Call the service layer
	order, err := h.orderService.CreateOrder(ctx, req.CustomerId, items, payments)
	if err != nil {
		log.Println("Failed to create order:", err)
		return nil, mapErrorToGRPCStatus(err).Err()
	}

	// Convert domain model to proto
	protoOrder := utils.OrderToProto(order)

	return &pb.CreateOrderResponse{
		Order: protoOrder,
	}, nil
}

func (h *OrderHandler) CreateInvoice(ctx context.Context, req *pb.CreateInvoiceRequest) (*pb.CreateInvoiceResponse, error) {
	// Convert proto items to DTO
	items := utils.ProtoInvoiceItemRequestsToDTO(req.Items)

	// Call service
	invoice, err := h.invoiceService.CreateInvoice(ctx, req.ShipmentId, req.OrderId, items)
	if err != nil {
		return &pb.CreateInvoiceResponse{
			Code:    "ERROR",
			Message: err.Error(),
		}, nil
	}

	// Convert domain model to proto
	protoInvoice := utils.InvoiceToProto(invoice)

	return &pb.CreateInvoiceResponse{
		Code:    "SUCCESS",
		Message: "Invoice created successfully",
		Invoice: protoInvoice,
	}, nil
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
