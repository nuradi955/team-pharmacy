package services

import (
	"errors"

	"team-pharmacy/internal/dto"
	"team-pharmacy/internal/models"
	"team-pharmacy/internal/repository"
)

type SubcategoryService interface {
	Create(req dto.SubcategoryCreate) (*models.Subcategory, error)
	GetByCategory(categoryID uint) ([]models.Subcategory, error)
	GetByID(id uint) (*models.Subcategory, error)
}

type subcategoryService struct {
	subRepo repository.SubcategoryRepository
	catRepo repository.CategoryRepository
}

func NewSubcategoryService(
	subRepo repository.SubcategoryRepository,
	catRepo repository.CategoryRepository,
) SubcategoryService {
	return &subcategoryService{
		subRepo: subRepo,
		catRepo: catRepo,
	}
}

func (s *subcategoryService) Create(req dto.SubcategoryCreate) (*models.Subcategory, error) {

_, err := s.catRepo.GetByID(req.CategoryID)
if err != nil {
	return nil, fmt.Errorf("category not found: %w", err)
}

	sub := &models.Subcategory{
		Name:       req.Name,
		CategoryID: req.CategoryID,
	}

	if err := s.subRepo.Create(sub); err != nil {
		return nil, err
	}

	return sub, nil
}

func (s *subcategoryService) GetByCategory(categoryID uint) ([]models.Subcategory, error) {
	return s.subRepo.ListByCategory(categoryID)
}

func (s *subcategoryService) GetByID(id uint) (*models.Subcategory, error) {
	return s.subRepo.GetByID(id)
}
