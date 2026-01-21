package repository

import (
	"team-pharmacy/internal/models"

	"gorm.io/gorm"
)

type CategoryRepository interface {
	List() ([]models.Category, error)
	Create(category *models.Category) error
	GetByID(id uint) (*models.Category, error)
	Update(category *models.Category) error
	Delete(id uint) error
}

type gormCategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &gormCategoryRepository{db: db}
}

func (r *gormCategoryRepository) List() ([]models.Category, error) {
	var categories []models.Category
	if err := r.db.Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *gormCategoryRepository) Create(category *models.Category) error {
	return r.db.Create(category).Error
}

func (r *gormCategoryRepository) GetByID(id uint) (*models.Category, error) {
	var category models.Category
	if err := r.db.First(&category, id).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *gormCategoryRepository) Update(category *models.Category) error {
	if category == nil {
		return nil
	}
	return r.db.Save(category).Error
}

func (r *gormCategoryRepository) Delete(id uint) error {
	return r.db.Delete(&models.Category{}, id).Error
}
