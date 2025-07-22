package repository

import (
	"context"
	"time"

	"github.com/best-microservice/order-service/internal/models"
	"github.com/google/uuid"

	"github.com/jmoiron/sqlx"
)

type OrderRepository struct {
	db *sqlx.DB
}

func NewOrderRepository(db *sqlx.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) CreateOrder(ctx context.Context, order *models.Order) error {
	query := `
		INSERT INTO users (id, user_id, item, total, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	order.ID = uuid.New().String()
	now := time.Now()

	err := r.db.QueryRowContext(ctx, query,
		order.ID, order.UserID, order.Items, order.Total, order.Status, now).Scan(&order.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *OrderRepository) GetOrderByID(ctx context.Context, id string) (*models.Order, error) {
	query := `
		SELECT *
		FROM users
		WHERE id = $1
	`

	var order models.Order
	err := r.db.GetContext(ctx, &order, query, id)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

type OrderItemRepository struct {
	db *sqlx.DB
}

func NewOrderItemRepository(db *sqlx.DB) *OrderItemRepository {
	return &OrderItemRepository{db: db}
}

func (r *OrderItemRepository) CreateOrderItem(ctx context.Context, orderItem *models.OrderItem) error {
	query := `
		INSERT INTO users (product_id, quantity, price)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	orderItem.ProductID = uuid.New().String()

	err := r.db.QueryRowContext(ctx, query,
		orderItem.ProductID, orderItem.Price, orderItem.Quantity).Scan(&orderItem.ProductID)
	if err != nil {
		return err
	}

	return nil
}

func (r *OrderItemRepository) GetOrderItemByID(ctx context.Context, id string) (*models.OrderItem, error) {
	query := `
		SELECT *
		FROM users
		WHERE id = $1
	`

	var orderItem models.OrderItem
	err := r.db.GetContext(ctx, &orderItem, query, id)
	if err != nil {
		return nil, err
	}

	return &orderItem, nil
}
