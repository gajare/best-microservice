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
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

func main() {
	// Initialize database connection
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

	// Run database migrations
	if err := runMigrations(db.DB); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
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

	// Gracefully stop gRPC server
	grpcServer.GracefulStop()

	// Close database connection with timeout
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second) //ctx replace by ctx
	defer cancel()
	if err := db.Close(); err != nil {
		log.Printf("error closing database connection: %v", err)
	}

	log.Println("Order service stopped")
}

// runMigrations executes database migrations
func runMigrations(db *sql.DB) error {
	// In a real application, you would use a proper migration tool like golang-migrate
	// This is a simplified version for demonstration

	// Check if tables exist
	var exists bool
	err := db.QueryRow(`
		SELECT EXISTS (
			SELECT FROM information_schema.tables 
			WHERE table_name = 'orders'
		)
	`).Scan(&exists)

	if err != nil {
		return fmt.Errorf("failed to check for existing tables: %w", err)
	}

	if exists {
		log.Println("Database tables already exist, skipping migrations")
		return nil
	}

	log.Println("Running database migrations...")

	// Create tables
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS orders (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			user_id UUID NOT NULL REFERENCES users(id),
			total DECIMAL(10,2) NOT NULL,
			status VARCHAR(50) NOT NULL DEFAULT 'pending',
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS order_items (
			order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
			product_id UUID NOT NULL REFERENCES products(id),
			quantity INTEGER NOT NULL,
			price DECIMAL(10,2) NOT NULL,
			PRIMARY KEY (order_id, product_id)
		);

		CREATE INDEX IF NOT EXISTS idx_orders_user_id ON orders(user_id);
	`)

	if err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}

	log.Println("Database migrations completed successfully")
	return nil
}
