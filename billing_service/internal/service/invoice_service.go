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
	//Validate order exists
	_, err := s.orderRepo.GetByID(ctx, orderId)
	if err != nil {
		return nil, fmt.Errorf("order not found: %w", err)
	}

	// Validate and process items
	var totalAmount float64
	var invoiceItems []model.InvoiceItem

	for _, itemReq := range itemRequest {
		if itemReq.Quantity <= 0 {
			return nil, fmt.Errorf("invalid quantity for item %s: %d", itemReq.Sku, itemReq.Quantity)
		}

		// Get item by SKU
		item, err := s.itemRepo.GetBySku(ctx, itemReq.Sku)
		if err != nil {
			return nil, fmt.Errorf("item with SKU %s not found: %w", itemReq.Sku, err)
		}

		itemTotal := item.Price * float64(itemReq.Quantity)
		totalAmount += itemTotal

		invoiceItems = append(invoiceItems, model.InvoiceItem{
			Quantity: itemReq.Quantity,
			ItemID:   item.ID,
		})
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
