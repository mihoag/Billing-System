package main

import (
	"log"
	"net"

	billing_handler "billing-system/billing_service/internal/handler"
	"billing-system/billing_service/internal/repository"
	"billing-system/billing_service/internal/service"
	"billing-system/billing_service/pkg/config"
	"billing-system/billing_service/pkg/db"
	billing_pb "billing-system/billing_service/proto"

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
	itemRepo := repository.NewItemRepository(gormDB)
	orderRepo := repository.NewOrderRepository(gormDB)

	// Initialize services
	orderService := service.NewOrderService(orderRepo, itemRepo)

	// Initialize HTTP handlers
	orderHandler := billing_handler.NewOrderHandler(orderService)

	// server's address
	address := "127.0.0.1:8082"
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Failed to listen on %s: %v", address, err)
	}

	// start gRPC server
	grpcServer := grpc.NewServer()
	billing_pb.RegisterBillingServiceServer(grpcServer, orderHandler)
	reflection.Register(grpcServer)

	if err = grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}

}
