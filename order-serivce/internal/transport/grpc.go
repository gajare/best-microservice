package transport

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/best-microservice/common/protos/order"
	"github.com/best-microservice/order-service/internal/models"
	"github.com/best-microservice/order-service/internal/service"
	"github.com/google/uuid"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OrderServer struct {
	order.UnimplementedOrderServiceServer
	service *service.OrderService
}

func NewOrderServer(service *service.OrderService) *OrderServer {
	return &OrderServer{service: service}
}

func (s *OrderServer) CreateOrder(ctx context.Context, req *order.CreateOrderRequest) (*order.OrderResponse, error) {
	// Validate request
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}
	if len(req.Items) == 0 {
		return nil, status.Error(codes.InvalidArgument, "at least one order item is required")
	}

	// Convert protobuf items to models
	var orderItems []models.OrderItem
	var total float64

	for _, item := range req.Items {
		if item.Quantity <= 0 {
			return nil, status.Errorf(codes.InvalidArgument, "invalid quantity for product %s", item.ProductId)
		}
		if item.Price <= 0 {
			return nil, status.Errorf(codes.InvalidArgument, "invalid price for product %s", item.ProductId)
		}

		orderItems = append(orderItems, models.OrderItem{
			ProductID: item.ProductId,
			Quantity:  int(item.Quantity),
			Price:     float64(item.Price),
		})

		total += float64(item.Price) * float64(item.Quantity)
	}

	// Create order model
	newOrder := &models.Order{
		ID:     uuid.New().String(),
		UserID: req.UserId,
		Items:  orderItems,
		Total:  total,
		Status: "pending",
	}

	// Call service layer
	err := s.service.CreateOrder(ctx, newOrder)
	if err != nil {
		if errors.Is(err, service.ErrProductNotFound) {
			return nil, status.Error(codes.NotFound, "product not found")
		}
		if errors.Is(err, service.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		if errors.Is(err, service.ErrInsufficientStock) {
			return nil, status.Error(codes.FailedPrecondition, "insufficient stock")
		}
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to create order: %v", err))
	}

	// Convert back to protobuf response
	return s.orderToResponse(newOrder), nil
}

func (s *OrderServer) GetOrder(ctx context.Context, req *order.GetOrderRequest) (*order.OrderResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "order id is required")
	}

	order, err := s.service.GetOrder(ctx, req.Id)
	if err != nil {
		if errors.Is(err, service.ErrOrderNotFound) {
			return nil, status.Error(codes.NotFound, "order not found")
		}
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to get order: %v", err))
	}

	return s.orderToResponse(order), nil
}

func (s *OrderServer) GetUserOrders(ctx context.Context, req *order.GetUserOrdersRequest) (*order.GetUserOrdersResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user id is required")
	}

	if req.Limit <= 0 {
		req.Limit = 10
	}
	if req.Offset < 0 {
		req.Offset = 0
	}

	orders, total, err := s.service.GetUserOrders(ctx, req.UserId, int(req.Limit), int(req.Offset))
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to get user orders: %v", err))
	}

	// Convert orders to protobuf response
	var pbOrders []*order.OrderResponse
	for _, o := range orders {
		pbOrders = append(pbOrders, s.orderToResponse(o))
	}

	return &order.GetUserOrdersResponse{
		Orders: pbOrders,
		Total:  int32(total),
	}, nil
}

func (s *OrderServer) orderToResponse(o *models.Order) *order.OrderResponse {
	var items []*order.OrderItem
	for _, item := range o.Items {
		items = append(items, &order.OrderItem{
			ProductId: item.ProductID,
			Quantity:  int32(item.Quantity),
			Price:     float32(item.Price),
		})
	}

	return &order.OrderResponse{
		Id:        o.ID,
		UserId:    o.UserID,
		Items:     items,
		Total:     float32(o.Total),
		Status:    o.Status,
		CreatedAt: o.CreatedAt.Format(time.RFC3339),
	}
}
