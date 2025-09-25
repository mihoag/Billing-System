package dto

import "billing-system/billing_service/internal/model"

// ItemRequest represents a request to include an item in an order or invoice
type ItemRequest struct {
	Sku      string `json:"skus"`
	Quantity int    `json:"quantity"`
}

// PaymentRequest represents a request to add a payment to an order
type PaymentRequest struct {
	Method model.PaymentMethod `json:"method"`
	Amount float64             `json:"amount"`
}
