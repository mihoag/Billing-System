package utils

import (
	"billing-system/billing_service/internal/dto"
	"billing-system/billing_service/internal/model"
	pb "billing-system/billing_service/proto"
	"time"
)

// Proto to DTO conversions
// ProtoItemRequestsToDTO converts protocol buffer item requests to DTO item requests
func ProtoItemRequestsToDTO(protoItems []*pb.ItemRequest) []dto.ItemRequest {
	if protoItems == nil {
		return nil
	}

	items := make([]dto.ItemRequest, len(protoItems))
	for i, item := range protoItems {
		items[i] = ProtoItemRequestToDTO(item)
	}
	return items
}

// ProtoItemRequestToDTO converts a single protocol buffer item request to a DTO item request
func ProtoItemRequestToDTO(protoItem *pb.ItemRequest) dto.ItemRequest {
	return dto.ItemRequest{
		Sku:      protoItem.Sku,
		Quantity: int(protoItem.Quantity),
	}
}

// ProtoPaymentRequestsToDTO converts protocol buffer payment requests to DTO payment requests
func ProtoPaymentRequestsToDTO(protoPayments []*pb.PaymentRequest) []dto.PaymentRequest {
	if protoPayments == nil {
		return nil
	}

	payments := make([]dto.PaymentRequest, len(protoPayments))
	for i, payment := range protoPayments {
		payments[i] = ProtoPaymentRequestToDTO(payment)
	}
	return payments
}

// ProtoPaymentRequestToDTO converts a single protocol buffer payment request to a DTO payment request
func ProtoPaymentRequestToDTO(protoPayment *pb.PaymentRequest) dto.PaymentRequest {
	return dto.PaymentRequest{
		Method: model.PaymentMethod(protoPayment.Method),
		Amount: protoPayment.Amount,
	}
}

// ProtoInvoiceItemRequestsToDTO converts protocol buffer invoice item requests to DTO invoice item requests
func ProtoInvoiceItemRequestsToDTO(protoItems []*pb.InvoiceItemRequest) []dto.InvoiceItemRequest {
	if protoItems == nil {
		return nil
	}

	items := make([]dto.InvoiceItemRequest, len(protoItems))
	for i, item := range protoItems {
		items[i] = ProtoInvoiceItemRequestToDTO(item)
	}
	return items
}

// ProtoInvoiceItemRequestToDTO converts a single protocol buffer invoice item request to a DTO invoice item request
func ProtoInvoiceItemRequestToDTO(protoItem *pb.InvoiceItemRequest) dto.InvoiceItemRequest {
	return dto.InvoiceItemRequest{
		Sku:      protoItem.Sku,
		Quantity: int(protoItem.Quantity),
	}
}

// Domain to Proto conversions

// OrderToProto converts a domain order model to a protocol buffer order
func OrderToProto(order *model.Order) *pb.Order {
	if order == nil {
		return nil
	}

	protoOrder := &pb.Order{
		Id:          order.ID,
		CustomerId:  order.CustomerID,
		TotalAmount: order.TotalAmount,
		Status:      OrderStatusToProto(order.Status),
		CreatedAt:   order.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   order.UpdatedAt.Format(time.RFC3339),
		Items:       OrderItemsToProto(order.Items),
		Payments:    PaymentsToProto(order.Payments),
	}

	return protoOrder
}

// OrderItemsToProto converts domain order items to protocol buffer order items
func OrderItemsToProto(items []model.OrderItem) []*pb.OrderItem {
	if items == nil {
		return nil
	}

	protoItems := make([]*pb.OrderItem, len(items))
	for i, item := range items {
		protoItems[i] = OrderItemToProto(&item)
	}
	return protoItems
}

// OrderItemToProto converts a single domain order item to a protocol buffer order item
func OrderItemToProto(item *model.OrderItem) *pb.OrderItem {
	if item == nil {
		return nil
	}

	return &pb.OrderItem{
		Id:       item.ID,
		OrderId:  item.OrderID,
		ItemId:   item.ItemID,
		Quantity: int32(item.Quantity),
	}
}

// PaymentsToProto converts domain payments to protocol buffer payments
func PaymentsToProto(payments []model.Payment) []*pb.Payment {
	if payments == nil {
		return nil
	}

	protoPayments := make([]*pb.Payment, len(payments))
	for i, payment := range payments {
		protoPayments[i] = PaymentToProto(&payment)
	}
	return protoPayments
}

// PaymentToProto converts a single domain payment to a protocol buffer payment
func PaymentToProto(payment *model.Payment) *pb.Payment {
	if payment == nil {
		return nil
	}

	return &pb.Payment{
		Id:      payment.ID,
		OrderId: payment.OrderID,
		Method:  string(payment.Method),
		Amount:  payment.Amount,
	}
}

// InvoiceToProto converts a domain invoice to a protocol buffer invoice
func InvoiceToProto(invoice *model.Invoice) *pb.Invoice {
	if invoice == nil {
		return nil
	}

	protoInvoice := &pb.Invoice{
		Id:          invoice.ID,
		ShipmentId:  invoice.ShipmentID,
		OrderId:     invoice.OrderID,
		TotalAmount: invoice.TotalAmount,
		CreatedAt:   invoice.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   invoice.UpdatedAt.Format(time.RFC3339),
		Items:       InvoiceItemsToProto(invoice.Items),
	}

	return protoInvoice
}

// InvoiceItemsToProto converts domain invoice items to protocol buffer invoice items
func InvoiceItemsToProto(items []model.InvoiceItem) []*pb.InvoiceItem {
	if items == nil {
		return nil
	}

	protoItems := make([]*pb.InvoiceItem, len(items))
	for i, item := range items {
		protoItems[i] = InvoiceItemToProto(&item)
	}
	return protoItems
}

// InvoiceItemToProto converts a single domain invoice item to a protocol buffer invoice item
func InvoiceItemToProto(item *model.InvoiceItem) *pb.InvoiceItem {
	if item == nil {
		return nil
	}

	return &pb.InvoiceItem{
		Id:        item.ID,
		InvoiceId: item.InvoiceID,
		ItemId:    item.ItemID,
		Quantity:  int32(item.Quantity),
	}
}

// OrderStatusToProto maps a domain OrderStatus to a proto OrderStatus
func OrderStatusToProto(status model.OrderStatus) pb.OrderStatus {
	switch status {
	case model.OrderPending:
		return pb.OrderStatus_PENDING
	case model.OrderSuccess:
		return pb.OrderStatus_SUCCESS
	case model.OrderFailed:
		return pb.OrderStatus_FAILED
	default:
		return pb.OrderStatus_PENDING
	}
}
