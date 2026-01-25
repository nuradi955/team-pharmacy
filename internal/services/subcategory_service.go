package services

import (
	"errors"
	"fmt"

	"team-pharmacy/internal/dto"
	"team-pharmacy/internal/models"
	"team-pharmacy/internal/repository"
)

var ErrCategoryNotFound = errors.New("category not found")

type SubcategoryService interface {
	Create(categoryID uint, req dto.SubcategoryCreateRequest) (*models.Subcategory, error)
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

func (s *subcategoryService) Create(categoryID uint, req dto.SubcategoryCreateRequest) (*models.Subcategory, error) {
	_, err := s.catRepo.GetByID(categoryID)
	if err != nil {
		return nil, fmt.Errorf("%w", ErrCategoryNotFound)
	}

	sub := &models.Subcategory{
		Name:       req.Name,
		CategoryID: categoryID,
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
