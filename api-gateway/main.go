package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/best-microservice/api-gateway/handlers"
)

func main() {
	// Initialize Gin router
	router := gin.Default()

	// Setup gRPC connections
	userConn, err := grpc.Dial(os.Getenv("USER_SERVICE_ADDR"),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect to user service: %v", err)
	}
	defer userConn.Close()

	productConn, err := grpc.Dial(os.Getenv("PRODUCT_SERVICE_ADDR"),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect to product service: %v", err)
	}
	defer productConn.Close()

	orderConn, err := grpc.Dial(os.Getenv("ORDER_SERVICE_ADDR"),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect to order service: %v", err)
	}
	defer orderConn.Close()

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userConn)
	productHandler := handlers.NewProductHandler(productConn)
	orderHandler := handlers.NewOrderHandler(orderConn)

	// Routes
	api := router.Group("/api/v1")
	{
		// User routes
		api.POST("/users", userHandler.CreateUser)
		api.GET("/users/:id", userHandler.GetUser)
		api.POST("/auth", userHandler.Authenticate)

		// Product routes
		api.POST("/products", productHandler.CreateProduct)
		api.GET("/products/:id", productHandler.GetProduct)
		api.GET("/products", productHandler.ListProducts)

		// Order routes
		api.POST("/orders", orderHandler.CreateOrder)
		api.GET("/orders/:id", orderHandler.GetOrder)
		api.GET("/users/:id/orders", orderHandler.GetUserOrders)
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Start server
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
	log.Println("Server exiting")
}
