package model

import "time"

// ShipmentStatus defines the status of a shipment
type ShipmentStatus string

const (
	Confirmed ShipmentStatus = "CONFIRMED"
	Failed    ShipmentStatus = "FAILED"
)

type Base struct {
	ID        int64      `json:"id" gorm:"primaryKey;autoIncrement"`
	CreatedAt time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}

// Shipment represents a shipment in the system
type Shipment struct {
	Base
	OrderID int64          `json:"order_id"`
	Items   []ShipmentItem `json:"items" gorm:"foreignKey:ShipmentID"`
	Status  ShipmentStatus `json:"status"`
}

// ShipmentItem represents an item in a shipment
type ShipmentItem struct {
	ShipmentID int64  `json:"shipment_id" gorm:"primaryKey"`
	Sku        string `json:"sku" gorm:"primaryKey"`
	Quantity   int    `json:"quantity"`
}
