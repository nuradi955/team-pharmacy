package transport

import (
	"net/http"
	"strconv"

	"team-pharmacy/internal/dto"
	"team-pharmacy/internal/services"

	"github.com/gin-gonic/gin"
)

type SubcategoryHandler struct {
	service services.SubcategoryService
}

func NewSubcategoryHandler(service services.SubcategoryService) *SubcategoryHandler {
	return &SubcategoryHandler{service: service}
}

func (h *SubcategoryHandler) RegisterRoutes(r *gin.Engine) {
	categories := r.Group("/categories")

	{
		categories.GET("/:id/subcategories", h.GetByCategory)
		categories.POST("/:id/subcategories", h.Create)
	}
}

func (h *SubcategoryHandler) Create(c *gin.Context) {
	var req dto.SubcategoryCreate

	// Берем category_id из URL
	categoryID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category id"})
		return
	}

	req.CategoryID = uint(categoryID)

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	subcategory, err := h.service.Create(req)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, subcategory)
}

func (h *SubcategoryHandler) GetByCategory(c *gin.Context) {
	categoryID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category id"})
		return
	}

	subcategories, err := h.service.GetByCategory(uint(categoryID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, subcategories)
}
