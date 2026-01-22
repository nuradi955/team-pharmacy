package services

import (
	"team-pharmacy/internal/dto"
	"team-pharmacy/internal/models"
	"team-pharmacy/internal/repository"
)

type CategoryService interface {
	CreateCategory(req dto.CategoryCreate) (*models.Category, error)
	GetList() ([]models.Category, error)
	GetByID(id uint) (*models.Category, error)
}

type categoryService struct {
	repo repository.CategoryRepository
}

func NewCategoryService(repo repository.CategoryRepository) CategoryService {
	return &categoryService{repo: repo}
}

func (s *categoryService) CreateCategory(req dto.CategoryCreate) (*models.Category, error) {
	category := &models.Category{Name: req.Name}

	if err := s.repo.Create(category); err != nil {
		return nil, err
	}

	return category, nil
}

func (s *categoryService) GetList() ([]models.Category, error) {
	return s.repo.List()
}

func (s *categoryService) GetByID(id uint) (*models.Category, error) {
	return s.repo.GetByID(id)
}
