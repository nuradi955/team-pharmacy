package repository

import (
	"team-pharmacy/internal/models"

	"gorm.io/gorm"
)

type SubcategoryRepository interface {
	ListByCategory(categoryID uint) ([]models.Subcategory, error)
	Create(subcategory *models.Subcategory) error
	GetByID(id uint) (*models.Subcategory, error)
	Update(subcategory *models.Subcategory) error
	Delete(id uint) error
}

type gormSubcategoryRepository struct {
	db *gorm.DB
}

func NewSubcategoryRepository(db *gorm.DB) SubcategoryRepository {
	return &gormSubcategoryRepository{db: db}
}

func (r *gormSubcategoryRepository) ListByCategory(categoryID uint) ([]models.Subcategory, error) {
	var subcategories []models.Subcategory

	if err := r.db.
		Where("category_id = ?", categoryID).
		Find(&subcategories).Error; err != nil {
		return nil, err
	}

	return subcategories, nil
}

func (r *gormSubcategoryRepository) Create(subcategory *models.Subcategory) error {
	return r.db.Create(subcategory).Error
}

func (r *gormSubcategoryRepository) GetByID(id uint) (*models.Subcategory, error) {
	var subcategory models.Subcategory

	if err := r.db.First(&subcategory, id).Error; err != nil {
		return nil, err
	}

	return &subcategory, nil
}

func (r *gormSubcategoryRepository) Update(subcategory *models.Subcategory) error {
	return r.db.Updates(subcategory).Error
}

func (r *gormSubcategoryRepository) Delete(id uint) error {
	return r.db.Delete(&models.Subcategory{}, id).Error
}
