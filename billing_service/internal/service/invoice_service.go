package service

import (
	"billing-system/billing_service/internal/dto"
	"billing-system/billing_service/internal/model"
	"billing-system/billing_service/internal/repository"
	"context"
	"fmt"
)

type InvoiceServiceImpl struct {
	invoiceRepo repository.InvoiceRepository
	orderRepo   repository.OrderRepository
	itemRepo    repository.ItemRepository
}

func NewInvoiceService(invoiceRepo repository.InvoiceRepository, orderRepo repository.OrderRepository, itemRepo repository.ItemRepository) InvoiceService {
	return &InvoiceServiceImpl{
		invoiceRepo: invoiceRepo,
		orderRepo:   orderRepo,
		itemRepo:    itemRepo,
	}
}

func (s *InvoiceServiceImpl) CreateInvoice(ctx context.Context, shipmentId int64, orderId int64, itemRequest []dto.InvoiceItemRequest) (*model.Invoice, error) {
	// Validate order exists and get order details
	order, err := s.orderRepo.GetByID(ctx, orderId)
	if err != nil {
		return nil, fmt.Errorf("order not found: %w", err)
	}

	// Create a map of items in the order with quantities
	orderItemMap := make(map[int64]int)
	skuToItemID := make(map[string]int64)
	for _, orderItem := range order.Items {
		orderItemMap[orderItem.ItemID] = orderItem.Quantity
		skuToItemID[orderItem.Item.Sku] = orderItem.ItemID
	}

	// Get existing invoices for this order
	existingInvoices, err := s.invoiceRepo.GetByOrderID(ctx, orderId)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve existing invoices: %w", err)
	}

	// Calculate consumed quantities per item ID across all invoices
	consumedQuantities := make(map[int64]int)
	for _, invoice := range existingInvoices {
		for _, item := range invoice.Items {
			consumedQuantities[item.ItemID] += item.Quantity
		}
	}

	// Validate and process items
	var totalAmount float64
	var invoiceItems []model.InvoiceItem
	requestedQuantities := make(map[int64]int)

	for _, itemReq := range itemRequest {
		if itemReq.Quantity <= 0 {
			return nil, fmt.Errorf("invalid quantity for item %s: %d", itemReq.Sku, itemReq.Quantity)
		}

		// Get item by SKU
		item, err := s.itemRepo.GetBySku(ctx, itemReq.Sku)
		if err != nil {
			return nil, fmt.Errorf("item with SKU %s not found: %w", itemReq.Sku, err)
		}

		// Check if item exists in original order
		orderQty, exists := orderItemMap[item.ID]
		if !exists {
			return nil, fmt.Errorf("item %s not found in original order", itemReq.Sku)
		}

		// Check if the total quantity exceeds the order quantity
		totalRequestedQty := itemReq.Quantity + consumedQuantities[item.ID]
		if totalRequestedQty > orderQty {
			return nil, fmt.Errorf(
				"requested quantity %d for item %s exceeds available quantity %d (consumed: %d, ordered: %d)",
				itemReq.Quantity,
				itemReq.Sku,
				orderQty-consumedQuantities[item.ID],
				consumedQuantities[item.ID],
				orderQty,
			)
		}

		// Track requested quantities for this invoice
		requestedQuantities[item.ID] += itemReq.Quantity

		itemTotal := item.Price * float64(itemReq.Quantity)
		totalAmount += itemTotal

		invoiceItems = append(invoiceItems, model.InvoiceItem{
			Quantity: itemReq.Quantity,
			ItemID:   item.ID,
		})
	}

	// Double-check for duplicates in the current request
	for itemID, qty := range requestedQuantities {
		if qty > orderItemMap[itemID] {
			return nil, fmt.Errorf("duplicate items in request exceed original order quantity")
		}
	}

	// Create invoice
	invoice := &model.Invoice{
		OrderID:     orderId,
		ShipmentID:  shipmentId,
		TotalAmount: totalAmount,
		Items:       invoiceItems,
	}

	if err := s.invoiceRepo.Create(ctx, invoice); err != nil {
		return nil, fmt.Errorf("failed to create invoice: %w", err)
	}

	return invoice, nil
}
