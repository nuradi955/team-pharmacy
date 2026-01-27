package transport

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"team-pharmacy/internal/dto"
	"team-pharmacy/internal/errs"
	"team-pharmacy/internal/services"

	"github.com/gin-gonic/gin"
)

type CartHandler struct {
	service services.CartService
	logger  *slog.Logger
}

func NewCartHandler(logger *slog.Logger, service services.CartService) *CartHandler {
	return &CartHandler{service: service, logger: logger.With("layer", "transport", "entity", "cart")}
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
		h.logger.Warn("invalid user ID", "param", c.Param("id"), "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("incoming request", "method", c.Request.Method, "user_id", userID)

	cart, err := h.service.GetCartWithItems(uint(userID))
	if err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {

			h.logger.Warn("user not found", "user_id", userID)

			c.JSON(http.StatusNotFound, err.Error())
			return
		}
		h.logger.Error("failed to get cart", "user_id", userID, "error", err)

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return

	}
	h.logger.Info("cart returned", "user_id", userID, "items_count", len(cart.Items))

	c.JSON(http.StatusOK, cart)

}

func (h *CartHandler) CreateItem(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		h.logger.Warn("invalid  userID", 
		"raw_value", c.Param("id"),  
		"error", err,
	)

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var req *dto.AddCartItemRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid request body", "user_id", userID, "error", err)

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("incoming request",
		"method", c.Request.Method,
		"path", c.FullPath(),
		"user_id", userID,
		"medicine_id",
		req.MedicineID,
		"quantity", req.Quantity,
	)

	cart, err := h.service.CreateItem(uint(userID), req)
	if err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {

			h.logger.Warn("user not found", "user_id", userID)

			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		h.logger.Error("failed to add item to cart",
			"user_id", userID,
			"medicine_id", req.MedicineID,
			"quantity", req.Quantity,
		)

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("item added to cart",
		"user_id", userID,
		"medicine_id", req.MedicineID,
		"quantity", req.Quantity,
	)

	c.JSON(http.StatusOK, cart)
}

func (h *CartHandler) UpdateItem(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {

		h.logger.Warn("invalid  userID",
			"raw_value", c.Param("id"),
			"user_id", userID, "error", err,
		)

		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	itemID, err := strconv.ParseUint(c.Param("item_id"), 10, 64)
	if err != nil {

		h.logger.Warn("invalid itemID",
			"raw_value", c.Param("item_id"),
			"user_id", userID, "error", err,
		)

		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid item id"})
		return
	}
	var req *dto.UpdateCartItemRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid request body",
			"user_id", userID,
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("incoming request",
		"user_id", userID,
		"item_id", itemID,
	)

	newItem, err := h.service.UpdateItem(uint(userID), uint(itemID), req)

	if err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			h.logger.Warn("user not found",
				"user_id", userID,
				"error", err,
			)
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		h.logger.Error("failed to update cart item",
			"user_id", userID,
			"item_id", itemID,
			"error", err,
		)

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("cart item updated",
		"user_id", userID,
		"item_id", itemID,
	)

	c.JSON(http.StatusOK, newItem)
}

func (h *CartHandler) DeleteItem(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {

		h.logger.Warn("invalid userID",
			"raw_value", c.Param("id"),
			"error", err,
		)

		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	itemID, err := strconv.ParseUint(c.Param("item_id"), 10, 64)
	if err != nil {

		h.logger.Warn("invalid itemID",
			"raw_value", c.Param("item_id"),
			"error", err,
		)

		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid item id"})
		return
	}

	h.logger.Info("incoming request",
		"user_id", userID,
		"item_id", itemID,
	)
	if err := h.service.DeleteItem(uint(userID), uint(itemID)); err != nil {
		if errors.Is(err, errs.ErrUserNotFound) || errors.Is(err, errs.ErrItemNotFound) {
			h.logger.Error("resource not found",
				"user_id", userID,
				"item_id", itemID,
				"error", err,
			)

			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		h.logger.Error("failed to delete cart item",
			"user_id", userID,
			"item_id", itemID,
			"error", err,
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	h.logger.Info("cart item deleted",
		"user_id", userID,
		"item_id", itemID,
	)

	c.Status(http.StatusOK)
}

func (h *CartHandler) ClearCart(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		h.logger.Warn("invalid userID",
			"raw_value", c.Param("id"),
			"error", err,
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("incoming request",
		"user_id", userID,
	)

	if err := h.service.ClearCart(uint(userID)); err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			h.logger.Warn("user not found",
				"user_id", userID)
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return

		}

		h.logger.Error("failed to clear cart",
			"user_id", userID,
			"error", err,
		)

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	h.logger.Info("cart cleared",
		"user_id", userID)

	c.Status(http.StatusOK)
}
