package billing

// CreateOrderRequest represents a request to create a new order
type CreateOrderRequest struct {
	CustomerID string           `json:"customer_id" binding:"required"`
	Items      []ItemRequest    `json:"items" binding:"required,dive"`
	Payments   []PaymentRequest `json:"payments" binding:"required,dive"`
}

// ItemRequest represents an item in a create order request
type ItemRequest struct {
	ItemID   int64   `json:"item_id" binding:"required"`
	Quantity int     `json:"quantity" binding:"required,min=1"`
	Price    float64 `json:"price" binding:"omitempty,min=0"`
}

// PaymentRequest represents a payment in a create order request
type PaymentRequest struct {
	Method string  `json:"method" binding:"required"`
	Amount float64 `json:"amount" binding:"required,min=0"`
}

// CreateOrderResponse represents a response from creating an order
type CreateOrderResponse struct {
	Order OrderResponse `json:"order"`
}

// OrderResponse represents an order in responses
type OrderResponse struct {
	ID          int64               `json:"id"`
	CustomerID  string              `json:"customer_id"`
	TotalAmount float64             `json:"total_amount"`
	Status      string              `json:"status"`
	Items       []OrderItemResponse `json:"items"`
	Payments    []PaymentResponse   `json:"payments"`
	CreatedAt   string              `json:"created_at"`
	UpdatedAt   string              `json:"updated_at"`
}

// OrderItemResponse represents an order item in responses
type OrderItemResponse struct {
	ID       int64 `json:"id"`
	OrderID  int64 `json:"order_id"`
	ItemID   int64 `json:"item_id"`
	Quantity int   `json:"quantity"`
}

// PaymentResponse represents a payment in responses
type PaymentResponse struct {
	ID      int64   `json:"id"`
	OrderID int64   `json:"order_id"`
	Method  string  `json:"method"`
	Amount  float64 `json:"amount"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}
