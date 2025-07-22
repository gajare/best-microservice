package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/best-microservice/order-service/internal/models"
	"github.com/best-microservice/order-service/internal/repository"
)

var (
	ErrOrderNotFound     = errors.New("order not found")
	ErrProductNotFound   = errors.New("product not found")
	ErrUserNotFound      = errors.New("user not found")
	ErrInsufficientStock = errors.New("insufficient stock")
)

type OrderService struct {
	repo *repository.OrderRepository
}

func NewOrderService(repo *repository.OrderRepository) *OrderService {
	return &OrderService{repo: repo}
}

func (s *OrderService) CreateOrder(ctx context.Context, order *models.Order) error {
	// Check if email already exists
	_, err := s.repo.GetOrderByID(ctx, order.ID)
	if err == nil {
		return errors.New("oder id  already exists")
	} else if err != sql.ErrNoRows {
		return err
	}

	return s.repo.CreateOrder(ctx, order)
}

func (s *OrderService) GetOder(ctx context.Context, id string) (*models.Order, error) {
	return s.repo.GetOrderByID(ctx, id)
}

type OrderItemService struct {
	repo *repository.OrderItemRepository
}

func NewOrderItemService(repo *repository.OrderItemRepository) *OrderItemService {
	return &OrderItemService{repo: repo}
}

func (s *OrderItemService) CreateOrderItem(ctx context.Context, orderItem *models.OrderItem) error {
	// Check if email already exists
	_, err := s.repo.GetOrderItemByID(ctx, orderItem.ProductID)
	if err == nil {
		return errors.New("oder id  already exists")
	} else if err != sql.ErrNoRows {
		return err
	}

	return s.repo.CreateOrderItem(ctx, orderItem)
}

func (s *OrderService) GetOrderByID(ctx context.Context, id string) (*models.Order, error) {
	return s.repo.GetOrderByID(ctx, id)
}

func (s *OrderService) GetOrder(ctx context.Context, id string) (*models.Order, error) {
	return nil, fmt.Errorf("GetOrder not implemented")
}
func (s *OrderService) GetUserOrders(ctx context.Context, userID string, limit, offset int) ([]*models.Order, int, error) {
	return nil, 0, fmt.Errorf("GetUserOrders not implemented")
}

// func (s *OrderService) GetOrder(ctx context.Context, id string) (*models.Order, error) {
// 	order, err := s.repo.GetOrderByID(ctx, id)
// 	if err != nil {
// 		if err == repository.ErrNotFound {
// 			return nil, ErrOrderNotFound
// 		}
// 		return nil, fmt.Errorf("failed to get order: %w", err)
// 	}

// 	// Enrich order items with product details
// 	for i := range order.Items {
// 		product, err := s.repo.GetProductByID(ctx, order.Items[i].ProductID)
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to get product details: %w", err)
// 		}
// 		order.Items[i].ProductName = product.Name
// 		order.Items[i].ProductDescription = product.Description
// 	}

// 	return order, nil
// }

// func (s *OrderService) GetUserOrders(ctx context.Context, userID string, limit, offset int) ([]*models.Order, int, error) {
// 	// Validate user exists
// 	_, err := s.userRepo.GetUserByID(ctx, userID)
// 	if err != nil {
// 		if err == repository.ErrNotFound {
// 			return nil, 0, ErrUserNotFound
// 		}
// 		return nil, 0, fmt.Errorf("failed to get user: %w", err)
// 	}

// 	orders, total, err := s.orderRepo.GetOrdersByUserID(ctx, userID, limit, offset)
// 	if err != nil {
// 		return nil, 0, fmt.Errorf("failed to get user orders: %w", err)
// 	}

// 	// Enrich order items with product details
// 	for _, order := range orders {
// 		for i := range order.Items {
// 			product, err := s.productRepo.GetProductByID(ctx, order.Items[i].ProductID)
// 			if err != nil {
// 				return nil, 0, fmt.Errorf("failed to get product details: %w", err)
// 			}
// 			order.Items[i].ProductName = product.Name
// 			order.Items[i].ProductDescription = product.Description
// 		}
// 	}

// 	return orders, total, nil
// }
