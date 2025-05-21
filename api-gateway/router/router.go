package router

import (
	"api-gateway/config"
	"api-gateway/handler"
	"github.com/gin-gonic/gin"
)

func SetupRouter(clients *config.Clients) *gin.Engine {
	r := gin.Default()

	users := r.Group("/users")
	{
		users.POST("/register", handler.RegisterUser(clients.UserClient))
		users.POST("/login", handler.AuthenticateUser(clients.UserClient))
		users.GET("/:id", handler.GetUserProfile(clients.UserClient))
		users.PUT("/:id", handler.UpdateUserProfile(clients.UserClient))
		users.DELETE("/:id", handler.DeleteUser(clients.UserClient))
		users.GET("/", handler.ListUsers(clients.UserClient))
	}

	products := r.Group("/products")
	{
		products.POST("", handler.CreateProduct(clients.InventoryClient))
		products.GET("/:id", handler.GetProduct(clients.InventoryClient))
		products.PUT("/:id", handler.UpdateProduct(clients.InventoryClient))
		products.DELETE("/:id", handler.DeleteProduct(clients.InventoryClient))
		products.GET("", handler.ListProducts(clients.InventoryClient))
	}

	orders := r.Group("/orders")
	{
		orders.POST("", handler.CreateOrder(clients.OrderClient))
		orders.GET("/:id", handler.GetOrderByID(clients.OrderClient))
		orders.PATCH("/:id/status", handler.UpdateOrderStatus(clients.OrderClient))
		orders.PUT("/:id", handler.UpdateOrder(clients.OrderClient))
		orders.DELETE("/:id", handler.DeleteOrderByID(clients.OrderClient))
		orders.GET("", handler.ListOrdersByUser(clients.OrderClient))
	}

	return r
}
