package repository

import (
	"context"
	"time"

	"github.com/best-microservice/product-service/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ProductRepository struct {
	db *sqlx.DB
}

func NewProductRepository(db *sqlx.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) CreateProduct(ctx context.Context, product *models.Product) error {
	query := `
		INSERT INTO products (id, name, description, price, stock, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`

	product.ID = uuid.New().String()
	now := time.Now()

	err := r.db.QueryRowContext(ctx, query,
		product.ID, product.Name, product.Description,
		product.Price, product.Stock, now, now).Scan(&product.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *ProductRepository) GetProductByID(ctx context.Context, id string) (*models.Product, error) {
	query := `
		SELECT id, name, description, price, stock, created_at
		FROM products
		WHERE id = $1
	`

	var product models.Product
	err := r.db.GetContext(ctx, &product, query, id)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *ProductRepository) ListProducts(ctx context.Context, limit, offset int) ([]*models.Product, int, error) {
	// Get total count
	var total int
	countQuery := `SELECT COUNT(*) FROM products`
	err := r.db.GetContext(ctx, &total, countQuery)
	if err != nil {
		return nil, 0, err
	}

	// Get products
	query := `
		SELECT id, name, description, price, stock, created_at
		FROM products
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	var products []*models.Product
	err = r.db.SelectContext(ctx, &products, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	return products, total, nil
}
