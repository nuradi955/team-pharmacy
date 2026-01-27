package services

import (
	"errors"
	"strings"
	"team-pharmacy/internal/dto"
	"team-pharmacy/internal/models"
	"team-pharmacy/internal/repository"
)

type MedicineService interface {
	Create(req dto.MedicineCreate) (*models.Medicine, error)
	GetAll() ([]models.Medicine, error)
	GetByID(id uint) (*models.Medicine, error)
	Update(req dto.MedicineUpdate, id uint) error
	Delete(id uint) error
}

type medicineService struct {
	MedicineRepo  repository.MedicineRepository
	CategoryRP    repository.CategoryRepository
	SubCategoryRP repository.SubcategoryRepository
}

func NewMedicineService(medicineRepo repository.MedicineRepository, categoryRepo repository.CategoryRepository, subcategoryRepo repository.SubcategoryRepository) MedicineService {
	return &medicineService{MedicineRepo: medicineRepo, CategoryRP: categoryRepo, SubCategoryRP: subcategoryRepo}
}

func (m *medicineService) Create(req dto.MedicineCreate) (*models.Medicine, error) {
	// Cheking Category and Subcategory
	_, err := m.CategoryRP.GetByID(*req.CategoryID)
	if err != nil {
		return nil, errors.New("Invalid Category")
	}
	_, err = m.SubCategoryRP.GetByID(req.SubcategoryID)
	if err != nil {
		return nil, errors.New("Invalid Subcategory")
	}
	// Cheking Name
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, errors.New("Write the Name")
	}
	// Cheking Price
	if req.Price <= 0 {
		return nil, errors.New("Isnt Correct Price")
	}

	// Create
	medicine := &models.Medicine{
		Name:                 name,
		Description:          req.Description,
		Price:                req.Price,
		StockQuantity:        req.StockQuantity,
		CategoryID:           req.CategoryID,
		SubcategoryID:        req.SubcategoryID,
		Manufacturer:         req.Manufacturer,
		PrescriptionRequired: req.PrescriptionRequired,
	}
	if err := m.MedicineRepo.Create(medicine); err != nil {
		return nil, err
	}
	return medicine, nil
}
func (m *medicineService) GetAll() ([]models.Medicine, error) {
	return m.MedicineRepo.GetAll()
}
func (m *medicineService) GetByID(id uint) (*models.Medicine, error) {
	if id == 0 {
		return nil, errors.New("Id can,t be zero")
	}
	medicine, err := m.MedicineRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return medicine, nil
}
func (m *medicineService) Update(req dto.MedicineUpdate, id uint) error {
	if id == 0 {
		return errors.New("Id cant be zero")
	}

	medicine, err := m.MedicineRepo.GetByID(id)
	if err != nil {
		return err
	}

	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if name == "" {
			return errors.New("name isnt correct")
		}
		medicine.Name = name
	}

	if req.Manufacturer != nil {
		manufacturer := strings.TrimSpace(*req.Manufacturer)
		if manufacturer == "" {
			return errors.New("Manufacturer isnt correct")
		}
		medicine.Manufacturer = manufacturer
	}

	if req.Price != nil {
		if *req.Price <= 0 {
			return errors.New("Isnt Correct Price")
		}
		medicine.Price = *req.Price
	}

	if req.StockQuantity != nil {
		medicine.StockQuantity = *req.StockQuantity
	}

	if req.CategoryID != nil {
		if _, err := m.CategoryRP.GetByID(*req.CategoryID); err != nil {
			return errors.New("CategoryID isnt Correct")
		}
		medicine.CategoryID = req.CategoryID
	}

	if req.SubcategoryID != nil {
		sub, err := m.SubCategoryRP.GetByID(*req.SubcategoryID)
		if err != nil {
			return errors.New("Subcategory isnt correct")
		}
		if medicine.CategoryID == nil {
			return errors.New("Error:CategoryID is nill")
		}
		if sub.CategoryID != *medicine.CategoryID {
			return errors.New("Subcategory dont have Correct category")
		}

		medicine.SubcategoryID = req.SubcategoryID

	}

	if req.Description != nil {
		medicine.Description = *req.Description
	}

	if req.PrescriptionRequired != nil {
		medicine.PrescriptionRequired = *req.PrescriptionRequired
	}

	if err := m.MedicineRepo.Update(medicine); err != nil {
		return err
	}
	return nil
}

func (m *medicineService) Delete(id uint) error {
	if id == 0 {
		return errors.New("Id cant be zero")
	}
	_, err := m.MedicineRepo.GetByID(id)
	if err != nil {
		return errors.New("this id isnt valid")
	}
	return m.MedicineRepo.Delete(id)
}
