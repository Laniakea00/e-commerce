package handler

import (
	"context"
	"net/http"
	"strconv"

	orderpb "github.com/Laniakea00/e-commerce/proto/order"
	"github.com/gin-gonic/gin"
)

func CreateOrder(client orderpb.OrderServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req orderpb.OrderRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Получаем user_id из контекста (уже валидный из JWT)
		userIDInterface, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		userID, ok := userIDInterface.(int32)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
			return
		}

		// Кладём userID в запрос к order-service
		req.UserId = userID

		resp, err := client.CreateOrder(context.Background(), &req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, resp)
	}
}

func GetOrderByID(client orderpb.OrderServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
			return
		}

		resp, err := client.GetOrderByID(context.Background(), &orderpb.OrderID{Id: int32(id)})
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, resp)
	}
}

func UpdateOrderStatus(client orderpb.OrderServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req orderpb.StatusUpdateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		resp, err := client.UpdateOrderStatus(context.Background(), &req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, resp)
	}
}

func ListOrdersByUser(client orderpb.OrderServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDInterface, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		userID, ok := userIDInterface.(int32)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
			return
		}

		res, err := client.ListOrdersByUser(context.Background(), &orderpb.UserID{
			Id: userID,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, res.Orders)
	}
}

func UpdateOrder(client orderpb.OrderServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req orderpb.Order
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		resp, err := client.UpdateOrder(context.Background(), &req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, resp)
	}
}

func DeleteOrderByID(client orderpb.OrderServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
			return
		}

		resp, err := client.DeleteOrderByID(context.Background(), &orderpb.OrderID{Id: int32(id)})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, resp)
	}
}
