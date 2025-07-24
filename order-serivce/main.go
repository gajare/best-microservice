package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/best-microservice/common/protos/order"
	"github.com/best-microservice/order-service/internal/repository"
	"github.com/best-microservice/order-service/internal/service"
	"github.com/best-microservice/order-service/internal/transport"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	// "github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Connect to the database
	db, err := sqlx.Connect("postgres", fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	))
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	// Load shared schema
	if err := loadSharedSchema(db.DB); err != nil {
		log.Fatalf("failed to load shared schema: %v", err)
	}

	// Initialize repository and service
	orderRepo := repository.NewOrderRepository(db)
	orderService := service.NewOrderService(orderRepo)

	// Create gRPC server
	grpcServer := grpc.NewServer()
	orderServer := transport.NewOrderServer(orderService)
	order.RegisterOrderServiceServer(grpcServer, orderServer)

	// Start gRPC server
	lis, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	go func() {
		log.Println("Starting order service on :50053")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down order service...")
	grpcServer.GracefulStop()

	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.Close(); err != nil {
		log.Printf("error closing database connection: %v", err)
	}
	log.Println("Order service stopped")
}

// loadSharedSchema loads the shared SQL schema into the DB (for dev/local use)
func loadSharedSchema(db *sql.DB) error {
	schema, err := os.ReadFile("../shared/schema.sql")
	if err != nil {
		return fmt.Errorf("failed to read schema.sql: %w", err)
	}

	_, err = db.Exec(string(schema))
	if err != nil {
		return fmt.Errorf("failed to execute schema.sql: %w", err)
	}

	log.Println("Shared schema applied successfully.")
	return nil
}
