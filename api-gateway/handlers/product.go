package handlers

import (
	"context"
	"net/http"
	"strconv"

	productpb "github.com/best-microservice/common/protos/product"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type ProductHandler struct {
	client productpb.ProductServiceClient
}

func NewProductHandler(conn *grpc.ClientConn) *ProductHandler {
	return &ProductHandler{
		client: productpb.NewProductServiceClient(conn),
	}
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req struct {
		Name        string  `json:"name" binding:"required"`
		Description string  `json:"description"`
		Price       float64 `json:"price" binding:"required,gt=0"`
		Stock       int     `json:"stock" binding:"gte=0"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	grpcReq := &productpb.CreateProductRequest{
		Name:        req.Name,
		Description: req.Description,
		Price:       float32(req.Price),
		Stock:       int32(req.Stock),
	}

	res, err := h.client.CreateProduct(context.Background(), grpcReq)
	if err != nil {
		st, _ := status.FromError(err)
		c.JSON(convertGRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":          res.Id,
		"name":        res.Name,
		"description": res.Description,
		"price":       res.Price,
		"stock":       res.Stock,
		"created_at":  res.CreatedAt,
	})
}

func (h *ProductHandler) GetProduct(c *gin.Context) {
	id := c.Param("id")

	res, err := h.client.GetProduct(context.Background(), &productpb.GetProductRequest{Id: id})
	if err != nil {
		st, _ := status.FromError(err)
		c.JSON(convertGRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":          res.Id,
		"name":        res.Name,
		"description": res.Description,
		"price":       res.Price,
		"stock":       res.Stock,
		"created_at":  res.CreatedAt,
	})
}

func (h *ProductHandler) ListProducts(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	res, err := h.client.ListProducts(context.Background(), &productpb.ListProductsRequest{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		st, _ := status.FromError(err)
		c.JSON(convertGRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}

	products := make([]gin.H, len(res.Products))
	for i, p := range res.Products {
		products[i] = gin.H{
			"id":          p.Id,
			"name":        p.Name,
			"description": p.Description,
			"price":       p.Price,
			"stock":       p.Stock,
			"created_at":  p.CreatedAt,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"products": products,
		"total":    res.Total,
		"limit":    limit,
		"offset":   offset,
	})
}
