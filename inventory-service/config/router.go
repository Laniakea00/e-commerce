package config

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"inventory-service/handler"
	"inventory-service/repository"
	"inventory-service/usecase"
)

func SetupRouter(db *sql.DB) *gin.Engine {
	r := gin.Default()

	productRepo := repository.NewProductRepository(db)
	productUsecase := usecase.NewProductUsecase(productRepo)
	productHandler := handler.NewProductHandler(productUsecase)

	products := r.Group("/products")
	{
		products.POST("", productHandler.CreateProduct)       // POST /products
		products.GET("/", productHandler.ListProducts)        // GET /products/
		products.GET("/:id", productHandler.GetProduct)       // GET /products?id=1
		products.PATCH("/:id", productHandler.UpdateProduct)  // PATCH /products/:id
		products.DELETE("/:id", productHandler.DeleteProduct) // DELETE /products/:id
	}

	return r
}
