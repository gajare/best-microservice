package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/best-microservice/product-service/internal/models"
	"github.com/best-microservice/product-service/internal/repository"
)

type ProductService struct {
	repo *repository.ProductRepository
}

func NewProductService(repo *repository.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (p *ProductService) CreatProduct(ctx context.Context, product *models.Product) error {
	_, err := p.repo.GetProductByID(ctx, product.ID)
	if err != nil {
		return errors.New("product already exist")
	} else if err != sql.ErrNoRows {
		return err
	}
	return p.repo.CreateProduct(ctx, product)
}

func (p *ProductService) GetProduct(ctx context.Context, id string) (*models.Product, error) {
	return p.repo.GetProductByID(ctx, id)
}
