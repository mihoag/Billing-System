package model

import (
	"time"
)

// OrderStatus defines the status of an order
type OrderStatus string

const (
	Pending OrderStatus = "PENDING"
	Success OrderStatus = "SUCCESS"
	Failed  OrderStatus = "FAILED"
)

// PaymentMethod defines the method of payment
type PaymentMethod string

const (
	COD   PaymentMethod = "COD"
	VNPAY PaymentMethod = "VN_PAY"
)

type Base struct {
	ID        int64      `json:"id" gorm:"primaryKey;autoIncrement"`
	CreatedAt time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}

// Order represents an order in the system
type Order struct {
	Base
	CustomerID  string      `json:"customer_id"`
	TotalAmount float64     `json:"total_amount"`
	Status      OrderStatus `json:"status"`
	Items       []OrderItem `json:"items,omitempty" gorm:"foreignKey:OrderID"`
	Payments    []Payment   `json:"payments,omitempty" gorm:"foreignKey:OrderID"`
	Invoices    []Invoice   `json:"invoices,omitempty" gorm:"foreignKey:OrderID"`
}

type Item struct {
	Base
	Name  string  `json:"name"`
	Sku   string  `json:"sku" gorm:"uniqueIndex"`
	Price float64 `json:"price"`
}

// OrderItem represents an item in an order
type OrderItem struct {
	Base
	OrderID  int64 `json:"order_id"`
	Quantity int   `json:"quantity"`
	ItemId   int64 `json:"item_id"`
}

// Payment represents a payment for an order
type Payment struct {
	Base
	OrderID int64         `json:"order_id"`
	Method  PaymentMethod `json:"method"`
	Amount  float64       `json:"amount"`
}

// Invoice represents an invoice for a shipment
type Invoice struct {
	Base
	OrderID     int64         `json:"order_id"`
	ShipmentID  int64         `json:"shipment_id" gorm:"uniqueIndex"`
	TotalAmount float64       `json:"total_amount"`
	Items       []InvoiceItem `json:"items" gorm:"foreignKey:InvoiceID"`
}

// InvoiceItem represents an item in an invoice
type InvoiceItem struct {
	Base
	InvoiceID int64 `json:"invoice_id"`
	Quantity  int   `json:"quantity"`
	ItemId    int64 `json:"item_id"`
}
