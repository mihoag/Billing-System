package service

import (
	"billing-system/shipment_service/client/billing"
	"billing-system/shipment_service/internal/dto"
	"billing-system/shipment_service/internal/model"
	"billing-system/shipment_service/internal/repository"
	"context"
	"fmt"
	"log"
)

type ShipmentServiceImpl struct {
	shipmentRepo  repository.ShipmentRepository
	billingClient *billing.BillingClient
}

func NewShipmentService(shipmentRepo repository.ShipmentRepository) ShipmentService {
	return &ShipmentServiceImpl{
		shipmentRepo:  shipmentRepo,
		billingClient: billing.NewBillingClient(),
	}
}

func (s *ShipmentServiceImpl) CreateShipment(ctx context.Context, orderID int64, items []dto.ShipmentItemRequest) (*model.Shipment, error) {
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

	// Create invoice in billing service
	invoiceItems := make([]billing.InvoiceItemRequest, len(shipment.Items))
	for i, item := range shipment.Items {
		invoiceItems[i] = billing.InvoiceItemRequest{
			Sku:      item.Sku,
			Quantity: int32(item.Quantity),
		}
	}

	invoiceReq := billing.CreateInvoiceRequest{
		ShipmentID: shipment.ID,
		OrderID:    shipment.OrderID,
		Items:      invoiceItems,
	}

	createInvoiceResponse, _ := s.billingClient.CreateInvoice(ctx, invoiceReq)
	if createInvoiceResponse.Code == "ERROR" {
		// Update shipment status to Failed
		shipment.Status = model.Failed
		if updateErr := s.shipmentRepo.Update(ctx, shipment); updateErr != nil {
			log.Printf("Failed to update shipment status: %v", updateErr)
		}

		return nil, fmt.Errorf("failed to create invoice: %s", createInvoiceResponse.Message)
	}

	return shipment, nil
}
