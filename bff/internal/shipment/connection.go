package shipment

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"billing-system/bff/config"
	shipmentPb "billing-system/shipment_service/proto"
)

type ShipmentConnectionAdapter struct {
	conn *grpc.ClientConn
}

func (shipmentConnectionAdapter *ShipmentConnectionAdapter) NewConnection() (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(
		// config.Service.BillingConnection.Address,
		config.Service.ShipmentConnection.Address, // Shipment service address
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}
	shipmentConnectionAdapter.conn = conn
	return conn, nil
}

func (shipmentConnectionAdapter *ShipmentConnectionAdapter) NewClient() (any, *grpc.ClientConn, error) {
	if shipmentConnectionAdapter.conn == nil {
		conn, err := shipmentConnectionAdapter.NewConnection()
		if err != nil {
			return nil, nil, err
		}
		shipmentConnectionAdapter.conn = conn
	}

	shipmentClient := shipmentPb.NewShipmentServiceClient(shipmentConnectionAdapter.conn)
	return shipmentClient, shipmentConnectionAdapter.conn, nil
}
