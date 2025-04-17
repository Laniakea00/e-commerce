package config

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"order-service/handler"
	"order-service/repository"
	"order-service/usecase"
)

func SetupRouter(db *sql.DB) *gin.Engine {
	r := gin.Default()
	orderRepo := repository.NewOrderRepository(db)
	orderUseCase := usecase.NewOrderUsecase(orderRepo)
	orderHandler := handler.NewOrderHandler(orderUseCase)

	orders := r.Group("/orders")
	{
		orders.POST("", orderHandler.CreateOrder)       // POST /orders
		orders.GET("/:id", orderHandler.GetOrder)       // GET /orders/:id
		orders.PATCH("/:id", orderHandler.UpdateStatus) // PATCH /orders/:id
		orders.PUT("/:id", orderHandler.UpdateOrder)    // PUT /orders/:id
		orders.DELETE("/:id", orderHandler.DeleteOrder) // DELETE /orders/:id
		orders.GET("", orderHandler.ListUserOrders)     // GET /orders?user_id=1
	}

	return r
}
