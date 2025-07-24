package handlers

import (
	"context"
	"net/http"
	"strconv"

	orderpb "github.com/best-microservice/common/protos/order"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type OrderHandler struct {
	client orderpb.OrderServiceClient
}

func NewOrderHandler(conn *grpc.ClientConn) *OrderHandler {
	return &OrderHandler{
		client: orderpb.NewOrderServiceClient(conn),
	}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req struct {
		UserID string `json:"user_id" binding:"required"`
		Items  []struct {
			ProductID string  `json:"product_id" binding:"required"`
			Quantity  int     `json:"quantity" binding:"required,gt=0"`
			Price     float64 `json:"price" binding:"required,gt=0"`
		} `json:"items" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert items to protobuf format
	items := make([]*orderpb.OrderItem, len(req.Items))
	for i, item := range req.Items {
		items[i] = &orderpb.OrderItem{
			ProductId: item.ProductID,
			Quantity:  int32(item.Quantity),
			Price:     float32(item.Price),
		}
	}

	grpcReq := &orderpb.CreateOrderRequest{
		UserId: req.UserID,
		Items:  items,
	}

	res, err := h.client.CreateOrder(context.Background(), grpcReq)
	if err != nil {
		st, _ := status.FromError(err)
		c.JSON(convertGRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}

	// Convert response items
	responseItems := make([]gin.H, len(res.Items))
	for i, item := range res.Items {
		responseItems[i] = gin.H{
			"product_id": item.ProductId,
			"quantity":   item.Quantity,
			"price":      item.Price,
		}
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":         res.Id,
		"user_id":    res.UserId,
		"items":      responseItems,
		"total":      res.Total,
		"status":     res.Status,
		"created_at": res.CreatedAt,
	})
}

func (h *OrderHandler) GetOrder(c *gin.Context) {
	id := c.Param("id")

	res, err := h.client.GetOrder(context.Background(), &orderpb.GetOrderRequest{Id: id})
	if err != nil {
		st, _ := status.FromError(err)
		c.JSON(convertGRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}

	// Convert response items
	items := make([]gin.H, len(res.Items))
	for i, item := range res.Items {
		items[i] = gin.H{
			"product_id": item.ProductId,
			"quantity":   item.Quantity,
			"price":      item.Price,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"id":         res.Id,
		"user_id":    res.UserId,
		"items":      items,
		"total":      res.Total,
		"status":     res.Status,
		"created_at": res.CreatedAt,
	})
}

func (h *OrderHandler) GetUserOrders(c *gin.Context) {
	userID := c.Param("id")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	res, err := h.client.GetUserOrders(context.Background(), &orderpb.GetUserOrdersRequest{
		UserId: userID,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		st, _ := status.FromError(err)
		c.JSON(convertGRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}

	orders := make([]gin.H, len(res.Orders))
	for i, order := range res.Orders {
		items := make([]gin.H, len(order.Items))
		for j, item := range order.Items {
			items[j] = gin.H{
				"product_id": item.ProductId,
				"quantity":   item.Quantity,
				"price":      item.Price,
			}
		}

		orders[i] = gin.H{
			"id":         order.Id,
			"user_id":    order.UserId,
			"items":      items,
			"total":      order.Total,
			"status":     order.Status,
			"created_at": order.CreatedAt,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"orders": orders,
		"total":  res.Total,
		"limit":  limit,
		"offset": offset,
	})
}
