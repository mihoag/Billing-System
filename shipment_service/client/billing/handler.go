package billing

import (
	billingPb "billing-system/billing_service/proto"
	"context"
	"fmt"
)

// BillingClient provides methods to interact with the billing service
type BillingClient struct {
	Connection *BillingConnectionAdapter
}

// NewBillingClient creates a new billing client
func NewBillingClient() *BillingClient {
	return &BillingClient{
		Connection: &BillingConnectionAdapter{},
	}
}

// CreateInvoice calls the billing service to create an invoice for a shipment
func (c *BillingClient) CreateInvoice(ctx context.Context, req CreateInvoiceRequest) (*CreateInvoiceResponse, error) {
	// Get billing service client
	clientInterface, _, err := c.Connection.NewClient()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to billing service: %w", err)
	}

	billingClient := clientInterface.(billingPb.BillingServiceClient)

	// Convert request to protobuf
	pbItems := make([]*billingPb.InvoiceItemRequest, len(req.Items))
	for i, item := range req.Items {
		pbItems[i] = &billingPb.InvoiceItemRequest{
			Sku:      item.Sku,
			Quantity: item.Quantity,
		}
	}

	pbRequest := &billingPb.CreateInvoiceRequest{
		ShipmentId: req.ShipmentID,
		OrderId:    req.OrderID,
		Items:      pbItems,
	}

	// Call billing service
	pbResponse, err := billingClient.CreateInvoice(ctx, pbRequest)
	if err != nil {
		return nil, fmt.Errorf("error calling billing service: %w", err)
	}

	// Convert response
	response := &CreateInvoiceResponse{
		Code:    pbResponse.Code,
		Message: pbResponse.Message,
		Invoice: pbResponse.Invoice,
	}

	return response, nil
}
