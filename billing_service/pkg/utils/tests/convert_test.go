package tests

import (
	"billing-system/billing_service/internal/dto"
	"billing-system/billing_service/internal/model"
	"billing-system/billing_service/pkg/utils"
	pb "billing-system/billing_service/proto"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestProtoItemRequestsToDTO(t *testing.T) {
	tests := []struct {
		name        string
		protoItems  []*pb.ItemRequest
		expectedDTO []dto.ItemRequest
	}{
		{
			name: "Success - Convert multiple proto items to DTO",
			protoItems: []*pb.ItemRequest{
				{Sku: "SKU-001", Quantity: 2, Price: 25.0},
				{Sku: "SKU-002", Quantity: 1, Price: 50.0},
			},
			expectedDTO: []dto.ItemRequest{
				{Sku: "SKU-001", Quantity: 2},
				{Sku: "SKU-002", Quantity: 1},
			},
		},
		{
			name:        "Success - Empty proto items",
			protoItems:  []*pb.ItemRequest{},
			expectedDTO: []dto.ItemRequest{},
		},
		{
			name:        "Success - Nil proto items",
			protoItems:  nil,
			expectedDTO: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.ProtoItemRequestsToDTO(tt.protoItems)
			
			if tt.protoItems == nil {
				assert.Nil(t, result)
				return
			}
			
			assert.Equal(t, len(tt.expectedDTO), len(result))
			
			for i, expected := range tt.expectedDTO {
				assert.Equal(t, expected.Sku, result[i].Sku)
				assert.Equal(t, expected.Quantity, result[i].Quantity)
			}
		})
	}
}

func TestOrderToProto(t *testing.T) {
	testTime := time.Now()
	
	tests := []struct {
		name        string
		order       *model.Order
		expectedNil bool
		checkFields func(t *testing.T, proto *pb.Order, order *model.Order)
	}{
		{
			name: "Success - Convert order to proto with items and payments",
			order: &model.Order{
				Base: model.Base{
					ID:        1,
					CreatedAt: testTime,
					UpdatedAt: testTime,
				},
				CustomerID:  "customer-123",
				TotalAmount: 100.0,
				Status:      model.OrderPending,
				Items: []model.OrderItem{
					{
						Base:     model.Base{ID: 1},
						OrderID:  1,
						ItemID:   100,
						Quantity: 2,
					},
				},
				Payments: []model.Payment{
					{
						Base:    model.Base{ID: 1},
						OrderID: 1,
						Method:  model.COD,
						Amount:  100.0,
					},
				},
			},
			expectedNil: false,
			checkFields: func(t *testing.T, proto *pb.Order, order *model.Order) {
				assert.Equal(t, order.ID, proto.Id)
				assert.Equal(t, order.CustomerID, proto.CustomerId)
				assert.Equal(t, order.TotalAmount, proto.TotalAmount)
				assert.Equal(t, pb.OrderStatus_PENDING, proto.Status)
				assert.Equal(t, order.CreatedAt.Format(time.RFC3339), proto.CreatedAt)
				assert.Equal(t, order.UpdatedAt.Format(time.RFC3339), proto.UpdatedAt)
				
				// Check items
				assert.Len(t, proto.Items, len(order.Items))
				for i, item := range order.Items {
					assert.Equal(t, item.ID, proto.Items[i].Id)
					assert.Equal(t, item.OrderID, proto.Items[i].OrderId)
					assert.Equal(t, item.ItemID, proto.Items[i].ItemId)
					assert.Equal(t, int32(item.Quantity), proto.Items[i].Quantity)
				}
				
				// Check payments
				assert.Len(t, proto.Payments, len(order.Payments))
				for i, payment := range order.Payments {
					assert.Equal(t, payment.ID, proto.Payments[i].Id)
					assert.Equal(t, payment.OrderID, proto.Payments[i].OrderId)
					assert.Equal(t, string(payment.Method), proto.Payments[i].Method)
					assert.Equal(t, payment.Amount, proto.Payments[i].Amount)
				}
			},
		},
		{
			name:        "Success - Nil order",
			order:       nil,
			expectedNil: true,
			checkFields: nil,
		},
		{
			name: "Success - Order with different status",
			order: &model.Order{
				Base: model.Base{
					ID:        2,
					CreatedAt: testTime,
					UpdatedAt: testTime,
				},
				CustomerID:  "customer-456",
				TotalAmount: 200.0,
				Status:      model.OrderSuccess,
			},
			expectedNil: false,
			checkFields: func(t *testing.T, proto *pb.Order, order *model.Order) {
				assert.Equal(t, pb.OrderStatus_SUCCESS, proto.Status)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.OrderToProto(tt.order)
			
			if tt.expectedNil {
				assert.Nil(t, result)
				return
			}
			
			assert.NotNil(t, result)
			if tt.checkFields != nil {
				tt.checkFields(t, result, tt.order)
			}
		})
	}
}

func TestOrderStatusToProto(t *testing.T) {
	tests := []struct {
		name           string
		status         model.OrderStatus
		expectedStatus pb.OrderStatus
	}{
		{
			name:           "Success - OrderPending to PENDING",
			status:         model.OrderPending,
			expectedStatus: pb.OrderStatus_PENDING,
		},
		{
			name:           "Success - OrderSuccess to SUCCESS",
			status:         model.OrderSuccess,
			expectedStatus: pb.OrderStatus_SUCCESS,
		},
		{
			name:           "Success - OrderFailed to FAILED",
			status:         model.OrderFailed,
			expectedStatus: pb.OrderStatus_FAILED,
		},
		{
			name:           "Success - Unknown status defaults to PENDING",
			status:         "UNKNOWN",
			expectedStatus: pb.OrderStatus_PENDING,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.OrderStatusToProto(tt.status)
			assert.Equal(t, tt.expectedStatus, result)
		})
	}
}