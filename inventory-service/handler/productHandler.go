package handler

import (
	"github.com/gin-gonic/gin"
	"inventory-service/domain"
	"inventory-service/usecase"
	"net/http"
	"strconv"
)

type ProductHandler struct {
	usecase usecase.ProductUsecase
}

func NewProductHandler(u usecase.ProductUsecase) *ProductHandler {
	return &ProductHandler{usecase: u}
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var p domain.Product
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.usecase.Create(p); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Product created"})
}

func (h *ProductHandler) GetProduct(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	product, err := h.usecase.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}
	c.JSON(http.StatusOK, product)
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var p domain.Product
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	p.ID = id
	if err := h.usecase.Update(p); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Product updated"})
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.usecase.DeleteByID(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Product deleted"})
}

func (h *ProductHandler) ListProducts(c *gin.Context) {
	products, err := h.usecase.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, products)
}
