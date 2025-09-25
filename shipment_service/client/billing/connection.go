package billing

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	billingPb "billing-system/billing_service/proto"
)

// BillingConnectionAdapter manages connections to the billing service
type BillingConnectionAdapter struct {
	conn *grpc.ClientConn
}

// NewConnection creates a new connection to the billing service
func (adapter *BillingConnectionAdapter) NewConnection() (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(
		//config.Service.BillingConnection.Address, // Billing service address
		"127.0.0.1:8082",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}
	adapter.conn = conn
	return conn, nil
}

// NewClient creates a new billing service client
func (adapter *BillingConnectionAdapter) NewClient() (any, *grpc.ClientConn, error) {
	if adapter.conn == nil {
		conn, err := adapter.NewConnection()
		if err != nil {
			return nil, nil, err
		}
		adapter.conn = conn
	}

	billingClient := billingPb.NewBillingServiceClient(adapter.conn)
	return billingClient, adapter.conn, nil
}
