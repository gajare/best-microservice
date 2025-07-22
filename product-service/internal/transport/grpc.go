package transport

import (
	"time"

	productpb "github.com/best-microservice/common/protos/product"
	"github.com/best-microservice/product-service/internal/models"
	"github.com/best-microservice/product-service/internal/service"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProductServer struct {
	productpb.UnimplementedProductServiceServer
	service *service.ProductService
}

func NewProductServer(service *service.ProductService) *ProductServer {
	return &ProductServer{service: service}
}

func (p *ProductServer) CreateProduct(ctx context.Context, req *productpb.CreateProductRequest) (*productpb.ProductResponse, error) {
	newProduct := &models.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       float64(req.Price),
		Stock:       int(req.Stock),
		CreatedAt:   time.Now(),
	}
	err := p.service.CreatProduct(ctx, newProduct)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "fail to creat product :%v", err)
	}
	return &productpb.ProductResponse{
		Id:          newProduct.ID,
		Name:        newProduct.Name,
		Description: newProduct.Description,
		Price:       float32(newProduct.Price),
		Stock:       int32(newProduct.Stock),
		// CreatedAt: newProduct.CreatedAt.String(),
	}, nil

}

func (p *ProductServer) GetProduct(ctx context.Context, req *productpb.GetProductRequest) (*productpb.ProductResponse, error) {
	product, err := p.service.GetProduct(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user not found :%v", err)
	}
	return &productpb.ProductResponse{
		Id:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Stock:       int32(product.Stock),
		Price:       float32(product.Price),
	}, nil
}
