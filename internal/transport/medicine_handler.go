package transport

import (
	"net/http"
	"strconv"
	"team-pharmacy/internal/dto"
	"team-pharmacy/internal/logger"
	"team-pharmacy/internal/services"

	"github.com/gin-gonic/gin"
)

type MedicineHandler struct {
	service services.MedicineService
}

func NewMedicineHandler(service services.MedicineService) *MedicineHandler {
	return &MedicineHandler{service: service}
}
func (m *MedicineHandler) RegisterRoutes(r *gin.Engine) {
	medicines := r.Group("/medicines")
	{
		medicines.GET("", m.GetAll)
		medicines.POST("", m.Create)
		medicines.GET("/:id", m.GetByID)
		medicines.PATCH("/:id", m.Update)
		medicines.DELETE("/:id", m.Delete)
	}
}
func (m *MedicineHandler) Create(ctx *gin.Context) {
	var req dto.MedicineCreate
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		logger.Log.Error(
			"Hadnler:Create medicine first step error",
			"error", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	medicine, err := m.service.Create(req)
	if err != nil {
		logger.Log.Error(
			"Handler:Create medicine second step error",
			"error", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	logger.Log.Info(
		"Handler:Create medicine correct",
		"medicine", medicine,
	)
	ctx.JSON(http.StatusCreated, medicine)
}

func (m *MedicineHandler) GetAll(ctx *gin.Context) {
	medicines, err := m.service.GetAll()
	if err != nil {
		
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, medicines)
}

func (m *MedicineHandler) GetByID(ctx *gin.Context) {
	id_check := ctx.Param("id")
	id, err := strconv.Atoi(id_check)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Is not Correct id"})
		return
	}
	if id <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id cant be < 0"})
		return
	}
	medicine, err := m.service.GetByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, medicine)
}

func (m *MedicineHandler) Update(ctx *gin.Context) {
	id_check := ctx.Param("id")
	id, err := strconv.Atoi(id_check)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Isnt correct id"})
		return
	}
	var medicine dto.MedicineUpdate
	if err := ctx.ShouldBindJSON(&medicine); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := m.service.Update(medicine, uint(id)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Update is Correct"})
}

func (m *MedicineHandler) Delete(ctx *gin.Context) {
	id_check := ctx.Param("id")
	id, err := strconv.Atoi(id_check)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Id isn,t correct"})
		return
	}
	if _, err := m.service.GetByID(uint(id)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "this medicine isnt exist"})
		return
	}
	if err := m.service.Delete(uint(id)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.Status(http.StatusNoContent)
}
