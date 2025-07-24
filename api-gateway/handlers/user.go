package handlers

import (
	"context"
	"net/http"

	userpb "github.com/best-microservice/common/protos/user"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type UserHandler struct {
	client userpb.UserServiceClient
}

func NewUserHandler(conn *grpc.ClientConn) *UserHandler {
	return &UserHandler{
		client: userpb.NewUserServiceClient(conn),
	}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=8"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	grpcReq := &userpb.CreateUserRequest{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	res, err := h.client.CreateUser(context.Background(), grpcReq)
	if err != nil {
		st, _ := status.FromError(err)
		c.JSON(convertGRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":         res.Id,
		"name":       res.Name,
		"email":      res.Email,
		"created_at": res.CreatedAt,
	})
}

func (h *UserHandler) GetUser(c *gin.Context) {
	id := c.Param("id")

	res, err := h.client.GetUser(context.Background(), &userpb.GetUserRequest{Id: id})
	if err != nil {
		st, _ := status.FromError(err)
		c.JSON(convertGRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":         res.Id,
		"name":       res.Name,
		"email":      res.Email,
		"created_at": res.CreatedAt,
	})
}

func (h *UserHandler) Authenticate(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.client.AuthenticateUser(context.Background(), &userpb.AuthRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		st, _ := status.FromError(err)
		c.JSON(convertGRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": res.Token,
		"user": gin.H{
			"id":         res.User.Id,
			"name":       res.User.Name,
			"email":      res.User.Email,
			"created_at": res.User.CreatedAt,
		},
	})
}

// func convertGRPCStatusToHTTP(code codes.Code) int {
// 	switch code {
// 	case codes.OK:
// 		return http.StatusOK
// 	case codes.InvalidArgument:
// 		return http.StatusBadRequest
// 	case codes.NotFound:
// 		return http.StatusNotFound
// 	case codes.AlreadyExists:
// 		return http.StatusConflict
// 	case codes.PermissionDenied:
// 		return http.StatusForbidden
// 	case codes.Unauthenticated:
// 		return http.StatusUnauthorized
// 	case codes.Internal:
// 		return http.StatusInternalServerError
// 	default:
// 		return http.StatusInternalServerError
// 	}
// }
