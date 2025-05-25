package router

import (
	"github.com/Laniakea00/e-commerce/api-gateway/config"
	"github.com/Laniakea00/e-commerce/api-gateway/handler"
	"github.com/Laniakea00/e-commerce/api-gateway/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/mattn/go-colorable"
	"net/http"
)

func SetupRouter(clients *config.Clients) *gin.Engine {
	gin.ForceConsoleColor()

	gin.DefaultWriter = colorable.NewColorableStdout()

	r := gin.Default()

	r.Use(cors.Default())

	r.LoadHTMLGlob("frontend/*.html")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	r.GET("/inventory", func(c *gin.Context) {
		c.HTML(http.StatusOK, "inventory.html", nil)
	})

	r.GET("/order", func(c *gin.Context) {
		c.HTML(http.StatusOK, "order.html", nil)
	})

	r.GET("/aggregated", middleware.AuthMiddleware(), handler.GetOrdersWithProductNames(clients.OrderClient, clients.InventoryClient))

	r.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", nil)
	})
	r.GET("/register", func(c *gin.Context) {
		c.HTML(http.StatusOK, "register.html", nil)
	})

	users := r.Group("/users")
	{
		users.POST("/register", handler.RegisterUser(clients.UserClient))
		users.POST("/login", handler.AuthenticateUser(clients.UserClient))
		users.GET("/:id", handler.GetUserProfile(clients.UserClient))
		users.PUT("/:id", handler.UpdateUserProfile(clients.UserClient))
		users.DELETE("/:id", handler.DeleteUser(clients.UserClient))
		users.GET("/", handler.ListUsers(clients.UserClient))
		users.GET("/verify", handler.VerifyUserEmail())
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
	orders.Use(middleware.AuthMiddleware())
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
