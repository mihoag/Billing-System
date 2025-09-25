package billing

// InvoiceItemRequest represents an item for invoice creation
type InvoiceItemRequest struct {
	Sku      string `json:"sku"`
	Quantity int32  `json:"quantity"`
}

// CreateInvoiceRequest represents the request to create an invoice
type CreateInvoiceRequest struct {
	ShipmentID int64                `json:"shipment_id"`
	OrderID    int64                `json:"order_id"`
	Items      []InvoiceItemRequest `json:"items"`
}

// CreateInvoiceResponse represents the response from creating an invoice
type CreateInvoiceResponse struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Invoice interface{} `json:"invoice,omitempty"`
}

// InvoiceData represents invoice data in a response
type InvoiceData struct {
	ID          int64   `json:"id"`
	ShipmentID  int64   `json:"shipment_id"`
	OrderID     int64   `json:"order_id"`
	TotalAmount float64 `json:"total_amount"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}
