package utils

import (
	"billing-system/shipment_service/internal/dto"
	"billing-system/shipment_service/internal/model"
	pb "billing-system/shipment_service/proto"
	"reflect"
	"testing"
	"time"
)

func TestConvertProtoItemsToDTO(t *testing.T) {
	// Test cases
	tests := []struct {
		name       string
		protoItems []*pb.ShipmentItemRequest
		want       []dto.ShipmentItemRequest
	}{
		{
			name:       "Empty slice",
			protoItems: []*pb.ShipmentItemRequest{},
			want:       []dto.ShipmentItemRequest{},
		},
		{
			name: "Single item",
			protoItems: []*pb.ShipmentItemRequest{
				{
					Sku:      "SKU123",
					Quantity: 5,
				},
			},
			want: []dto.ShipmentItemRequest{
				{
					Sku:      "SKU123",
					Quantity: 5,
				},
			},
		},
		{
			name: "Multiple items",
			protoItems: []*pb.ShipmentItemRequest{
				{
					Sku:      "SKU123",
					Quantity: 5,
				},
				{
					Sku:      "SKU456",
					Quantity: 10,
				},
				{
					Sku:      "SKU789",
					Quantity: 3,
				},
			},
			want: []dto.ShipmentItemRequest{
				{
					Sku:      "SKU123",
					Quantity: 5,
				},
				{
					Sku:      "SKU456",
					Quantity: 10,
				},
				{
					Sku:      "SKU789",
					Quantity: 3,
				},
			},
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConvertProtoItemsToDTO(tt.protoItems)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertProtoItemsToDTO() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConvertShipmentToProtoData(t *testing.T) {
	// Create a fixed timestamp for testing
	testTime := time.Date(2023, 9, 15, 12, 30, 0, 0, time.UTC)

	// Test cases
	tests := []struct {
		name     string
		shipment *model.Shipment
		want     *pb.ShipmentData
	}{
		{
			name: "Basic shipment",
			shipment: &model.Shipment{
				Base: model.Base{
					ID:        123,
					CreatedAt: testTime,
					UpdatedAt: testTime,
				},
				OrderID: 456,
				Status:  model.Confirmed,
				Items: []model.ShipmentItem{
					{
						ShipmentID: 123,
						Sku:        "SKU123",
						Quantity:   5,
					},
					{
						ShipmentID: 123,
						Sku:        "SKU456",
						Quantity:   10,
					},
				},
			},
			want: &pb.ShipmentData{
				ShipmentId: 123,
				OrderId:    456,
				Status:     string(model.Confirmed),
				CreatedAt:  testTime.Format(time.RFC3339),
				Items: []*pb.ShipmentItem{
					{
						Sku:      "SKU123",
						Quantity: 5,
					},
					{
						Sku:      "SKU456",
						Quantity: 10,
					},
				},
			},
		},
		{
			name: "Shipment with no items",
			shipment: &model.Shipment{
				Base: model.Base{
					ID:        123,
					CreatedAt: testTime,
					UpdatedAt: testTime,
				},
				OrderID: 456,
				Status:  model.Confirmed,
				Items:   []model.ShipmentItem{},
			},
			want: &pb.ShipmentData{
				ShipmentId: 123,
				OrderId:    456,
				Status:     string(model.Confirmed),
				CreatedAt:  testTime.Format(time.RFC3339),
				Items:      []*pb.ShipmentItem{},
			},
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConvertShipmentToProtoData(tt.shipment)

			// Check top-level fields
			if got.ShipmentId != tt.want.ShipmentId {
				t.Errorf("ShipmentId = %v, want %v", got.ShipmentId, tt.want.ShipmentId)
			}
			if got.OrderId != tt.want.OrderId {
				t.Errorf("OrderId = %v, want %v", got.OrderId, tt.want.OrderId)
			}
			if got.Status != tt.want.Status {
				t.Errorf("Status = %v, want %v", got.Status, tt.want.Status)
			}
			if got.CreatedAt != tt.want.CreatedAt {
				t.Errorf("CreatedAt = %v, want %v", got.CreatedAt, tt.want.CreatedAt)
			}

			// Check items
			if len(got.Items) != len(tt.want.Items) {
				t.Errorf("len(Items) = %v, want %v", len(got.Items), len(tt.want.Items))
			} else {
				for i, item := range got.Items {
					wantItem := tt.want.Items[i]
					if item.Sku != wantItem.Sku || item.Quantity != wantItem.Quantity {
						t.Errorf("Item[%d] = {Sku: %v, Quantity: %v}, want {Sku: %v, Quantity: %v}",
							i, item.Sku, item.Quantity, wantItem.Sku, wantItem.Quantity)
					}
				}
			}
		})
	}
}
