package services

import (
	"errors"

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
	if req.Name == "" {
		return nil, errors.New("category name must not be empty")
	}

	category := &models.Category{Name: req.Name}
	err := s.repo.Create(category)
	return category, err
}

func (s *categoryService) GetList() ([]models.Category, error) {
	return s.repo.List()
}

func (s *categoryService) GetByID(id uint) (*models.Category, error) {
	return s.repo.GetByID(id)
}
