package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"order-service/domain"
	"order-service/usecase"
)

type OrderHandler struct {
	uc usecase.OrderUsecase
}

func NewOrderHandler(uc usecase.OrderUsecase) *OrderHandler {
	return &OrderHandler{uc: uc}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var o domain.Order
	if err := c.ShouldBindJSON(&o); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.uc.Create(&o); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusCreated)
}

func (h *OrderHandler) GetOrder(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	order, err := h.uc.Get(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}
	c.JSON(http.StatusOK, order)
}

func (h *OrderHandler) UpdateStatus(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var data struct {
		Status string `json:"status"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := h.uc.UpdateStatus(id, data.Status)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}
	c.Status(http.StatusOK)
}

func (h *OrderHandler) ListUserOrders(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Query("user_id"))
	orders, err := h.uc.ListByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, orders)
}

func (h *OrderHandler) UpdateOrder(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id"})
		return
	}
	var o domain.Order
	if err := c.ShouldBindJSON(&o); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	o.ID = id
	err = h.uc.UpdateOrder(&o)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

func (h *OrderHandler) DeleteOrder(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	err := h.uc.DeleteByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}
	c.Status(http.StatusOK)
}
