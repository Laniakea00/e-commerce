package handler

import (
	"context"
	"net/http"
	"strconv"

	inventorypb "github.com/Laniakea00/e-commerce/proto/inventory"
	"github.com/gin-gonic/gin"
)

func CreateProduct(client inventorypb.InventoryServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req inventorypb.Product
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		resp, err := client.CreateProduct(context.Background(), &req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, resp)
	}
}

func GetProduct(client inventorypb.InventoryServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
			return
		}

		resp, err := client.GetProductByID(context.Background(), &inventorypb.ProductID{Id: int32(id)})
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, resp)
	}
}

func UpdateProduct(client inventorypb.InventoryServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req inventorypb.Product
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		resp, err := client.UpdateProduct(context.Background(), &req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, resp)
	}
}

func DeleteProduct(client inventorypb.InventoryServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
			return
		}

		resp, err := client.DeleteProductByID(context.Background(), &inventorypb.ProductID{Id: int32(id)})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, resp)
	}
}

func ListProducts(client inventorypb.InventoryServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		resp, err := client.ListProducts(context.Background(), &inventorypb.Empty{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, resp.Products)
	}
}
