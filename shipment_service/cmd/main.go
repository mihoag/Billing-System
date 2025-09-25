package main

import (
	"fmt"
	"log"
	"net"

	"billing-system/shipment_service/config"
	shipment_handler "billing-system/shipment_service/internal/handler"
	"billing-system/shipment_service/internal/repository"
	"billing-system/shipment_service/internal/service"
	"billing-system/shipment_service/pkg/db"

	shipment_pb "billing-system/shipment_service/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// load configuration from config.yaml
	// err := config.LoadConfig()
	// if err != nil {
	// 	log.Fatalf("Failed to get config: %v", err)
	// }

	gormDB, err := db.NewDatabase(&config.Service)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Run migrations
	if err := db.MigrateDB(gormDB); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize repositories
	shipmentRepo := repository.NewShipmentRepository(gormDB)

	// Initialize services
	shipmentService := service.NewShipmentService(shipmentRepo)

	// Initialize  handlers
	shipmentHandler := shipment_handler.NewShipmentHandler(shipmentService)

	// server's address
	address := fmt.Sprintf("%s:%s", config.Service.GRPCServer.Host, config.Service.GRPCServer.Port)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Failed to listen on %s: %v", address, err)
	}

	// start gRPC server
	grpcServer := grpc.NewServer()
	shipment_pb.RegisterShipmentServiceServer(grpcServer, shipmentHandler)
	reflection.Register(grpcServer)

	if err = grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}

}
