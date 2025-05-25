// handler/aggregated.go
package handler

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/Laniakea00/e-commerce/api-gateway/types"

	inventorypb "github.com/Laniakea00/e-commerce/proto/inventory"
	orderpb "github.com/Laniakea00/e-commerce/proto/order"
	"github.com/gin-gonic/gin"
)

func GetOrdersWithProductNames(
	orderClient orderpb.OrderServiceClient,
	inventoryClient inventorypb.InventoryServiceClient,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDStr := c.Query("user_id")
		if userIDStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
			return
		}

		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user_id must be a number"})
			return
		}

		ordersResp, err := orderClient.ListOrdersByUser(context.Background(), &orderpb.UserID{Id: int32(userID)})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get orders: " + err.Error()})
			return
		}

		var result []types.AggregatedOrder

		for _, o := range ordersResp.Orders {
			var items []types.OrderItem

			for _, item := range o.Items {
				productResp, err := inventoryClient.GetProductByID(context.Background(), &inventorypb.ProductID{Id: item.ProductId})
				if err != nil {
					log.Printf("product not found for ID %d: %v", item.ProductId, err)
					continue
				}

				items = append(items, types.OrderItem{
					ProductID:   item.ProductId,
					ProductName: productResp.Name,
					Quantity:    item.Quantity,
				})
			}

			result = append(result, types.AggregatedOrder{
				OrderID: o.Id,
				UserID:  o.UserId,
				Status:  o.Status,
				Items:   items,
			})
		}

		c.JSON(http.StatusOK, result)
	}
}
