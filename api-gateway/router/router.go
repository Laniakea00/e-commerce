package router

import (
	"api-gateway/proxy"
	"fmt"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Inventory Service
	r.Any("/products/*any", func(c *gin.Context) {
		target := fmt.Sprintf("http://localhost:8081/products%s", c.Param("any"))
		if c.Param("any") == "" || c.Param("any") == "/" {
			target = "http://localhost:8081/products"
		}
		proxy.Forward(c, target)
	})

	// Order Service
	r.Any("/orders/*any", func(c *gin.Context) {
		target := fmt.Sprintf("http://localhost:8082/orders%s", c.Param("any"))
		if c.Param("any") == "" || c.Param("any") == "/" {
			target = "http://localhost:8082/orders"
		}
		proxy.Forward(c, target)
	})

	return r
}
