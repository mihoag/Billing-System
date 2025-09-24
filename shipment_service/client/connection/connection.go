package connection

import (
	"google.golang.org/grpc"
)

type Connection interface {
	NewConnection() (*grpc.ClientConn, error)
	NewClient() (any, *grpc.ClientConn, error)
}
