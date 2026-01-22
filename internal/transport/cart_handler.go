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

type CartHandler struct {
	service services.CartService
}

func NewCartHandler(service services.CartService) *CartHandler {
	return &CartHandler{service: service}
}

func (h *CartHandler) RegisterRoutes(r *gin.Engine) {
	cart := r.Group("/users/:id/cart")
	{
		cart.GET("", h.GetCart)
		cart.POST("/items", h.CreateItem)
		cart.PATCH("/items/:item_id", h.UpdateItem)
		cart.DELETE("/items/:item_id", h.DeleteItem)
		cart.DELETE("", h.ClearCart)

	}

}

func (h *CartHandler) GetCart(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cart, err := h.service.GetCartWithItems(uint(userID))
	if err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, err.Error())
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cart)

}

func (h *CartHandler) CreateItem(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var req *dto.AddCartItemRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_ , err = h.service.CreateItem(uint(userID), req)
	if err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	cartWithItems, err := h.service.GetCartWithItems(uint(userID))
	if err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cartWithItems)
}

func (h *CartHandler) UpdateItem(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	itemID, err := strconv.ParseUint(c.Param("item_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid item id"})
		return
	}
	var req *dto.UpdateCartItemRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newItem, err := h.service.UpdateItem(uint(userID), uint(itemID), req)

	if err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, newItem)
}

func (h *CartHandler) DeleteItem(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	itemID, err := strconv.ParseUint(c.Param("item_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid item id"})
		return
	}

	if err := h.service.DeleteItem(uint(userID), uint(itemID)); err != nil {
		if errors.Is(err, errs.ErrUserNotFound) || errors.Is(err, errs.ErrItemNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

func (h *CartHandler) ClearCart(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.ClearCart(uint(userID)); err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
	}
	c.Status(http.StatusOK)
}
