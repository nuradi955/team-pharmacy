package transport

import (
	"net/http"
	"strconv"
	"team-pharmacy/internal/dto"
	"team-pharmacy/internal/services"

	"github.com/gin-gonic/gin"
)

type ReviewHandler struct {
	service services.ReviewService
}

func NewReviewHandler(service services.ReviewService) *ReviewHandler {
	return &ReviewHandler{service: service}
}

func (r *ReviewHandler) RegisterRoutes(g *gin.Engine) {
	reviews := g.Group("/reviews")
	{
		reviews.GET("/users/:id", r.GetAllByUser)
		reviews.GET("/medicines/:id", r.GetAllByMedicine)
		reviews.POST("", r.Create)
		reviews.GET("/:id", r.GetByID)
		reviews.PATCH("/:id", r.Update)
		reviews.DELETE("/:id", r.Delete)
	}
}

func (r *ReviewHandler) Create(ctx *gin.Context) {
	var req dto.ReviewCreate
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	review, err := r.service.Create(req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, review)
}
func (r *ReviewHandler) GetAllByUser(ctx *gin.Context) {
	id_check := ctx.Param("id")
	id, err := strconv.Atoi(id_check)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Isnt correct value user_id"})
		return
	}
	usersReviews, err := r.service.GetAllByUser(uint(id))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, usersReviews)
}

func (r *ReviewHandler) GetAllByMedicine(ctx *gin.Context) {
	id_check := ctx.Param("id")
	id, err := strconv.Atoi(id_check)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Isnt correct value medicine_id"})
		return
	}
	medicinesReviews, err := r.service.GetAllByMedicine(uint(id))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, medicinesReviews)
}

func (r *ReviewHandler) Update(ctx *gin.Context) {
	id_check := ctx.Param("id")
	id, err := strconv.Atoi(id_check)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Id isnt correct value"})
		return
	}
	var updateReviews dto.ReviewUpdate
	err = ctx.ShouldBindJSON(&updateReviews)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := r.service.Update(updateReviews, uint(id)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.Status(http.StatusOK)
}

func (r *ReviewHandler) Delete(ctx *gin.Context) {
	id_check := ctx.Param("id")
	id, err := strconv.Atoi(id_check)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Id isnt correct value"})
		return
	}
	if _, err := r.service.GetByID(uint(id)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "this review isnt exist"})
		return
	}
	if err := r.service.Delete(uint(id)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.Status(http.StatusNoContent)
}

func (r *ReviewHandler) GetByID(ctx *gin.Context) {
	id_check := ctx.Param("id")
	id, err := strconv.Atoi(id_check)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id isnt correct"})
		return
	}
	review, err := r.service.GetByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, review)
}
