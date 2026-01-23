package transport

import (
	"errors"
	"net/http"
	"strconv"
	"team-pharmacy/internal/dto"
	"team-pharmacy/internal/errs"
	"team-pharmacy/internal/services"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderService services.OrderService
	userService  services.UserService
	cartService  services.CartService
}

func NewOrderHandler(orderService services.OrderService, userService services.UserService, cartService services.CartService) *OrderHandler {
	return &OrderHandler{orderService: orderService, userService: userService, cartService: cartService}
}

func (h *OrderHandler) RegisterRoutes(r *gin.Engine) {
	order := r.Group("/orders/:id")
	{
		order.GET("", h.GetOrder)
		order.PATCH("/status", h.UpdateStatus)

	}
	user := r.Group("/users/:id")
	{
		user.POST("/orders", h.CreateOrder)
		user.GET("/orders", h.GetAllOrdersUser)
	}

}

func (h *OrderHandler) CreateOrder(c *gin.Context) {

	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var req dto.OrderCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := h.orderService.CreateOrder(uint(userID), &req)
	if err != nil {
		if errors.Is(err, errs.ErrCartNotFound) || errors.Is(err, errs.ErrCartIsEmpty) {
			c.JSON(http.StatusNotFound, gin.H{"error": "cart not found or is empty"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}
	c.JSON(http.StatusCreated, order)
}

func (h *OrderHandler) GetOrder(c *gin.Context) {
	orderID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := h.orderService.GetByID(uint(orderID))
	if err != nil {
		if errors.Is(err, errs.ErrOrderNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}
	c.JSON(http.StatusOK, order)
}

func (h *OrderHandler) GetAllOrdersUser(c *gin.Context) {

	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	orders, err := h.orderService.GetListOrders(uint(userID))
	if err != nil {
		if errors.Is(err, errs.ErrOrderNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "orders not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}
	c.JSON(http.StatusOK, orders)
}

func (h *OrderHandler) UpdateStatus(c *gin.Context) {
	orderID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var req dto.OrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.orderService.UpdateOrder(uint(orderID), &req); err != nil {
		if errors.Is(err, errs.ErrInvalidStatus) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "status error"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}
	c.Status(http.StatusOK)
}
