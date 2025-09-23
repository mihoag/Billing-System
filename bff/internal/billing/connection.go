package billing

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	billingPb "billing-system/billing_service/proto"
)

type BillingConnectionAdapter struct {
	conn *grpc.ClientConn
}

func (billingConnectionAdapter *BillingConnectionAdapter) NewConnection() (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(
		// config.Service.BillingConnection.Address,
		"127.0.0.1:8082",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}
	billingConnectionAdapter.conn = conn
	return conn, nil
}

func (billingConnectionAdapter *BillingConnectionAdapter) NewClient() (any, *grpc.ClientConn, error) {
	if billingConnectionAdapter.conn == nil {
		conn, err := billingConnectionAdapter.NewConnection()
		if err != nil {
			return nil, nil, err
		}
		billingConnectionAdapter.conn = conn
	}

	billingClient := billingPb.NewBillingServiceClient(billingConnectionAdapter.conn)
	return billingClient, billingConnectionAdapter.conn, nil
}
