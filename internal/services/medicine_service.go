package services

import (
	"team-pharmacy/internal/dto"
	"team-pharmacy/internal/models"
)

type MedicineService interface {
	Create(req dto.MedicineCreate) (*models.Medicine, error)
	GetAll() ([]models.Medicine, error)
	GetByID(id uint) (*models.Medicine, error)
	Update(req dto.MedicineUpdate, id uint) error
	Delete(id uint) error
}
