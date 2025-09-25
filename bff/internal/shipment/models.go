package shipment

// Request and response types
type ShipmentItemRequest struct {
	Sku      string `json:"sku"`
	Quantity int32  `json:"quantity"`
}

type CreateShipmentRequest struct {
	OrderID int64                 `json:"order_id"`
	Items   []ShipmentItemRequest `json:"items"`
}

type ShipmentResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
