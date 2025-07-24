package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/best-microservice/common/protos/user"
	"github.com/best-microservice/user-service/internal/repository"
	"github.com/best-microservice/user-service/internal/service"
	"github.com/best-microservice/user-service/internal/transport"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

func main() {
	// Database connection
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
	defer db.Close()

	// Initialize repository and service
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)

	// gRPC server
	grpcServer := grpc.NewServer()
	userServer := transport.NewUserServer(userService)
	user.RegisterUserServiceServer(grpcServer, userServer)

	// Start server
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	go func() {
		log.Println("Starting user service on :50051")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down user service...")
	grpcServer.GracefulStop()

	_, cancel := context.WithTimeout(context.Background(), 5*time.Second) //_ replace by ctx
	defer cancel()
	db.DB.Close()

	log.Println("User service stopped")
}
